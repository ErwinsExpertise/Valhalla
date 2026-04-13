package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Hucaru/Valhalla/nx"
	"github.com/Hucaru/gonx"
)

type candidate struct {
	File string
	Line int
	Kind string
	ID   int
	Text string
}

type missing struct {
	File string
	Line int
	Kind string
	ID   int
	Text string
}

var (
	itemCallRe    = regexp.MustCompile(`(?:itemCount|giveItem|removeItemsByID|takeItem|GetItem|CreateItemFromID|createPerfectItemFromID)\((\d{4,9})`)
	rewardRe      = regexp.MustCompile(`rewardID\s*=\s*(\d{4,9})`)
	storageRe     = regexp.MustCompile(`sendStorage\((\d{7,9})\)`)
	warpRe        = regexp.MustCompile(`(?:warp|plr\.warp)\([^\d]*(\d{7,9})`)
	pairItemRe    = regexp.MustCompile(`\[(\d{7,9})\s*,`)
	listItemRe    = regexp.MustCompile(`\b(\d{7,9})\b`)
	constAssignRe = regexp.MustCompile(`=\s*(\d{1,9})\b`)
	numberRe      = regexp.MustCompile(`\b(\d{4,9})\b`)
)

func main() {
	root := `C:\Users\daddy\desktop\MapleDev\Valhalla`
	v48 := filepath.Join(root, "..", "v48", "wz", "nx")

	nx.LoadFile(v48)

	faceSet := loadNameSet(filepath.Join(v48, "Character.nx"), "/Face")
	hairSet := loadNameSet(filepath.Join(v48, "Character.nx"), "/Hair")
	npcSet := loadNameSet(filepath.Join(v48, "Npc.nx"), "")

	var candidates []candidate
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" && ext != ".js" {
			return nil
		}
		candidates = append(candidates, scanFile(path)...)
		return nil
	})

	seen := make(map[string]bool)
	var missingList []missing
	for _, c := range candidates {
		key := fmt.Sprintf("%s:%d:%s:%d", c.File, c.Line, c.Kind, c.ID)
		if seen[key] {
			continue
		}
		seen[key] = true

		if exists(c, faceSet, hairSet, npcSet) {
			continue
		}

		missingList = append(missingList, missing{File: c.File, Line: c.Line, Kind: c.Kind, ID: c.ID, Text: strings.TrimSpace(c.Text)})
	}

	sort.Slice(missingList, func(i, j int) bool {
		if missingList[i].File != missingList[j].File {
			return missingList[i].File < missingList[j].File
		}
		if missingList[i].Line != missingList[j].Line {
			return missingList[i].Line < missingList[j].Line
		}
		if missingList[i].Kind != missingList[j].Kind {
			return missingList[i].Kind < missingList[j].Kind
		}
		return missingList[i].ID < missingList[j].ID
	})

	if len(missingList) == 0 {
		fmt.Println("No missing hard-coded NX IDs found in scanned source patterns.")
		return
	}

	for _, m := range missingList {
		fmt.Printf("%s:%d [%s] %d :: %s\n", strings.TrimPrefix(m.File, root+`\`), m.Line, m.Kind, m.ID, m.Text)
	}
}

func scanFile(path string) []candidate {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var out []candidate
	s := bufio.NewScanner(f)
	lineNo := 0
	for s.Scan() {
		lineNo++
		line := s.Text()
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "//") {
			continue
		}

		norm := filepath.ToSlash(path)
		switch {
		case strings.Contains(norm, "/constant/skill/class.go"):
			for _, m := range constAssignRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "skill", atoi(m[1]), line)
			}
		case strings.Contains(norm, "/constant/skill/mob.go"):
			for _, m := range constAssignRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "mob_skill", atoi(m[1]), line)
			}
		case strings.Contains(norm, "/constant/constants.go"):
			if strings.Contains(line, "Map") {
				for _, m := range constAssignRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "map", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "Mob") {
				for _, m := range constAssignRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "mob", atoi(m[1]), line)
				}
			}
		case strings.Contains(norm, "/login/handlers.go"):
			if strings.Contains(line, "allowedEyes") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "face", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "allowedHair") && !strings.Contains(line, "Colour") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "hair", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "allowedBottom") || strings.Contains(line, "allowedTop") || strings.Contains(line, "allowedShoes") || strings.Contains(line, "allowedWeapons") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
			if lineNo >= 430 && lineNo <= 438 {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
		case strings.Contains(norm, "/channel/commands.go"):
			if strings.Contains(line, "equips :=") || strings.Contains(line, "etc :=") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "= 2010007") || strings.Contains(line, "= 9200000") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "npc", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "= 5100001") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "mob", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "= 1010000") || strings.Contains(line, "= 60000") || strings.Contains(line, "= 104000000") || strings.Contains(line, "= 100000000") || strings.Contains(line, "= 103000000") || strings.Contains(line, "= 102000000") || strings.Contains(line, "= 101000000") || strings.Contains(line, "= 105040300") || strings.Contains(line, "= 180000000") || strings.Contains(line, "= 200000000") || strings.Contains(line, "= 211000000") || strings.Contains(line, "= 220000000") || strings.Contains(line, "= 221000000") || strings.Contains(line, "= 230000000") || strings.Contains(line, "= 105090900") || strings.Contains(line, "= 200000301") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "map", atoi(m[1]), line)
				}
			}
		case strings.HasSuffix(norm, ".js"):
			for _, m := range itemCallRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
			}
			for _, m := range rewardRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
			}
			for _, m := range storageRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "npc", atoi(m[1]), line)
			}
			for _, m := range warpRe.FindAllStringSubmatch(line, -1) {
				addCandidate(&out, path, lineNo, "map", atoi(m[1]), line)
			}
			if strings.Contains(line, "var items = [[") || strings.Contains(line, "const items = [[") {
				for _, m := range pairItemRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "var item = [") || strings.Contains(line, "var Prizes = [") {
				for _, m := range listItemRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
			if strings.Contains(line, "const ITEM_") || strings.Contains(line, "regCoupon") || strings.Contains(line, "vipCoupon") {
				for _, m := range numberRe.FindAllStringSubmatch(line, -1) {
					addCandidate(&out, path, lineNo, "item", atoi(m[1]), line)
				}
			}
		}
	}

	return out
}

func addCandidate(out *[]candidate, file string, line int, kind string, id int, text string) {
	if id <= 0 {
		return
	}
	*out = append(*out, candidate{File: file, Line: line, Kind: kind, ID: id, Text: text})
}

func exists(c candidate, faces, hairs, npcs map[int]bool) bool {
	switch c.Kind {
	case "item":
		_, err := nx.GetItem(int32(c.ID))
		return err == nil
	case "map":
		_, err := nx.GetMap(int32(c.ID))
		return err == nil
	case "mob":
		_, err := nx.GetMob(int32(c.ID))
		return err == nil
	case "quest":
		_, err := nx.GetQuest(int16(c.ID))
		return err == nil
	case "skill":
		_, err := nx.GetPlayerSkill(int32(c.ID))
		return err == nil
	case "mob_skill":
		_, err := nx.GetMobSkill(byte(c.ID))
		return err == nil
	case "npc":
		return npcs[c.ID]
	case "face":
		return faces[c.ID]
	case "hair":
		return hairs[c.ID]
	default:
		return true
	}
}

func loadNameSet(file string, path string) map[int]bool {
	set := make(map[int]bool)
	nodes, text, _, _, err := gonx.Parse(file)
	if err != nil {
		return set
	}

	var root *gonx.Node
	if path == "" {
		root = &nodes[0]
	} else if !gonx.FindNode(path, nodes, text, func(n *gonx.Node) { root = n }) {
		return set
	}

	for i := uint32(0); i < uint32(root.ChildCount); i++ {
		name := strings.TrimSuffix(text[nodes[root.ChildID+i].NameID], filepath.Ext(text[nodes[root.ChildID+i].NameID]))
		name = strings.TrimLeft(name, "0")
		if name == "" {
			name = "0"
		}
		if id, err := strconv.Atoi(name); err == nil {
			set[id] = true
		}
	}

	return set
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

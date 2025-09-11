package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Hucaru/gonx"
)

type questInfo struct {
	ID   int16
	Name string
}

func main() {
	var (
		nxPath     string
		scriptsDir string
		reportPath string
		dryRun     bool
		verbose    bool
	)
	flag.StringVar(&nxPath, "nx", "Data.nx", "Path to the NX file (e.g., Data.nx)")
	flag.StringVar(&scriptsDir, "scripts", "scripts/npc", "Path to the npc scripts directory")
	flag.StringVar(&reportPath, "report", "", "Optional path to write the report")
	flag.BoolVar(&dryRun, "dry-run", true, "If true, do not delete; only print what would be deleted")
	flag.BoolVar(&verbose, "v", true, "Verbose logging")
	flag.Parse()

	log.SetFlags(0)

	// Accept positional scripts dir too
	extras := flag.Args()
	if len(extras) == 1 {
		scriptsDir = extras[0]
	} else if len(extras) >= 2 {
		if strings.EqualFold(extras[0], "scripts") {
			scriptsDir = extras[1]
		} else if scriptsDir == "" || scriptsDir == "scripts/npc" {
			scriptsDir = extras[len(extras)-1]
		}
	}

	// Verify scripts directory exists and show absolute path for clarity
	scriptsDir = filepath.Clean(scriptsDir)
	absScriptsDir, _ := filepath.Abs(scriptsDir)
	if _, err := os.Stat(scriptsDir); err != nil {
		cwd, _ := os.Getwd()
		log.Printf("Failed to find scripts directory: %s", scriptsDir)
		log.Printf("Resolved absolute path: %s", absScriptsDir)
		log.Printf("Working directory: %s", cwd)
		log.Printf("Usage examples:")
		log.Printf("  go run quest-check.go -nx ..\\Data.nx -scripts .\\npc -report npc-quests.txt -dry-run")
		log.Printf("  go run quest-check.go -nx ..\\Data.nx .\\npc -report npc-quests.txt -dry-run")
		log.Fatalf("Error: %v", err)
	}

	if verbose {
		log.Printf("Loading NX from: %s", nxPath)
	}
	nodes, textLookup, _, _, err := gonx.Parse(nxPath)
	if err != nil {
		log.Fatalf("Failed to parse NX: %v", err)
	}

	if verbose {
		log.Printf("Extracting quest names and NPC associations from NX...")
	}
	questNames := extractQuestNames(nodes, textLookup)
	npcToQuestIDs := extractQuestNPCs(nodes, textLookup)

	// Build NPC -> []questInfo (attach names, uniq)
	npcToQuests := make(map[int32][]questInfo, len(npcToQuestIDs))
	for npcID, qIDs := range npcToQuestIDs {
		seen := make(map[int16]struct{}, len(qIDs))
		for _, qid := range qIDs {
			if _, ok := seen[qid]; ok {
				continue
			}
			seen[qid] = struct{}{}
			npcToQuests[npcID] = append(npcToQuests[npcID], questInfo{
				ID:   qid,
				Name: questNames[qid],
			})
		}
	}

	// Walk scripts and decide deletions:
	// delete only if:
	//  - file name is <npc-id>.js AND
	//  - npc-id has quests in NX AND
	//  - file contains the word "quest" (case-insensitive)
	if verbose {
		log.Printf("Scanning scripts directory: %s", absScriptsDir)
	}
	var toDelete []string
	err = filepath.WalkDir(scriptsDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".js") {
			return nil
		}

		base := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		id64, perr := strconv.ParseInt(base, 10, 32)
		if perr != nil {
			if verbose {
				log.Printf("Skipping non-numeric script file: %s", d.Name())
			}
			return nil
		}
		npcID := int32(id64)

		// Only consider scripts whose NPC is referenced by quests
		if _, referencedByQuests := npcToQuests[npcID]; !referencedByQuests {
			return nil
		}

		hasQuestWord, rerr := fileContainsQuestWord(path)
		if rerr != nil && verbose {
			log.Printf("Warning: failed reading %s: %v", path, rerr)
		}
		if hasQuestWord {
			toDelete = append(toDelete, path)
			if verbose && dryRun {
				log.Printf("[match] %s contains 'quest' and NPC %d has quest references", path, npcID)
			}
		} else if verbose {
			log.Printf("[skip] %s: NPC %d has quest references but file lacks 'quest'", path, npcID)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk scripts dir: %v", err)
	}

	// Report is always printed, regardless of deletions
	report := buildReport(npcToQuests)
	fmt.Print(report)

	// Optionally write report to file
	if reportPath != "" {
		if werr := os.WriteFile(reportPath, []byte(report), 0o644); werr != nil {
			log.Printf("Failed to write report to %s: %v", reportPath, werr)
		} else if verbose {
			log.Printf("Report written to: %s", reportPath)
		}
	}

	// Delete files or dry-run
	if len(toDelete) == 0 {
		log.Println("No NPC scripts matched both: 'has associated quests' AND contains 'quest'.")
		return
	}

	if dryRun {
		log.Printf("[DRY-RUN] Would delete %d files:", len(toDelete))
		for _, p := range toDelete {
			log.Println(" -", p)
		}
		return
	}

	var delCount, delErrs int
	for _, p := range toDelete {
		if err := os.Remove(p); err != nil {
			delErrs++
			log.Printf("Failed to delete %s: %v", p, err)
			continue
		}
		delCount++
		if verbose {
			log.Println("Deleted:", p)
		}
	}

	if delErrs > 0 {
		log.Printf("Done with errors. Deleted %d, %d failed.", delCount, delErrs)
	} else {
		log.Printf("Done. Deleted %d files.", delCount)
	}
}

func buildReport(npcToQuests map[int32][]questInfo) string {
	if len(npcToQuests) == 0 {
		return "No NPCs associated with quests were found.\n"
	}
	// Sort NPC IDs for stable output
	npcIDs := make([]int, 0, len(npcToQuests))
	for id := range npcToQuests {
		npcIDs = append(npcIDs, int(id))
	}
	sort.Ints(npcIDs)

	var b strings.Builder
	for _, id := range npcIDs {
		qs := npcToQuests[int32(id)]
		// Sort quests by ID
		sort.Slice(qs, func(i, j int) bool { return qs[i].ID < qs[j].ID })
		b.WriteString(fmt.Sprintf("%d:\n", id))
		for _, q := range qs {
			name := strings.TrimSpace(q.Name)
			if name == "" {
				name = "(unnamed quest)"
			}
			b.WriteString(fmt.Sprintf(" %d - %s\n", q.ID, name))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// extractQuestNames traverses /Quest/QuestInfo.img and returns questID -> name.
func extractQuestNames(nodes []gonx.Node, text []string) map[int16]string {
	out := make(map[int16]string)
	const root = "/Quest/QuestInfo.img"
	found := gonx.FindNode(root, nodes, text, func(n *gonx.Node) {
		for i := uint32(0); i < uint32(n.ChildCount); i++ {
			dir := nodes[n.ChildID+i]
			raw := text[dir.NameID]
			name := strings.TrimSuffix(raw, filepath.Ext(raw))
			qid64, err := strconv.ParseInt(name, 10, 16)
			if err != nil {
				continue
			}
			var qname string
			for j := uint32(0); j < uint32(dir.ChildCount); j++ {
				ch := nodes[dir.ChildID+j]
				key := text[ch.NameID]
				if key == "name" && ch.ChildCount == 0 {
					qname = text[gonx.DataToUint32(ch.Data)]
					break
				}
			}
			out[int16(qid64)] = qname
		}
	})
	if !found {
		log.Printf("Warning: could not find %s in NX.", root)
	}
	return out
}

// extractQuestNPCs traverses /Quest/Check.img and returns NPCID -> []questID
// by reading "npc" fields in both start ("0") and complete ("1") phases.
func extractQuestNPCs(nodes []gonx.Node, text []string) map[int32][]int16 {
	out := make(map[int32][]int16)
	const root = "/Quest/Check.img"
	found := gonx.FindNode(root, nodes, text, func(n *gonx.Node) {
		for i := uint32(0); i < uint32(n.ChildCount); i++ {
			dir := nodes[n.ChildID+i]
			raw := text[dir.NameID]
			name := strings.TrimSuffix(raw, filepath.Ext(raw))
			qid64, err := strconv.ParseInt(name, 10, 16)
			if err != nil {
				continue
			}
			qid := int16(qid64)

			for j := uint32(0); j < uint32(dir.ChildCount); j++ {
				phaseDir := nodes[dir.ChildID+j]
				for k := uint32(0); k < uint32(phaseDir.ChildCount); k++ {
					entry := nodes[phaseDir.ChildID+k]
					if text[entry.NameID] == "npc" && entry.ChildCount == 0 {
						npcID := gonx.DataToInt32(entry.Data)
						if npcID != 0 {
							out[npcID] = append(out[npcID], qid)
						}
					}
				}
			}
		}
	})
	if !found {
		log.Printf("Warning: could not find %s in NX.", root)
	}
	return out
}

// fileContainsQuestWord returns true if the file content contains the word "quest"
// (case-insensitive substring match).
func fileContainsQuestWord(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	s := strings.ToLower(string(data))
	return strings.Contains(s, "quest"), nil
}

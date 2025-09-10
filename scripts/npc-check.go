package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Hucaru/gonx"
)

func main() {
	var (
		nxPath     string
		scriptsDir string
		dryRun     bool
		verbose    bool
	)
	flag.StringVar(&nxPath, "nx", "Data.nx", "Path to the NX file (e.g., Data.nx)")
	flag.StringVar(&scriptsDir, "scripts", "scripts/npc", "Path to the npc scripts directory")
	flag.BoolVar(&dryRun, "dry-run", true, "If true, do not delete; only print what would be deleted")
	flag.BoolVar(&verbose, "v", true, "Verbose logging")
	flag.Parse()

	log.SetFlags(0)

	if verbose {
		log.Printf("Loading NX from: %s", nxPath)
	}
	nodes, textLookup, _, _, err := gonx.Parse(nxPath)
	if err != nil {
		log.Fatalf("Failed to parse NX: %v", err)
	}

	// Build a set of all node names for O(1) existence checks.
	// We only need names, not hierarchy.
	nxNames := make(map[string]struct{}, len(nodes))
	for i := range nodes {
		name := strings.ToLower(textLookup[nodes[i].NameID])
		if name != "" {
			nxNames[name] = struct{}{}
		}
	}

	if verbose {
		log.Printf("Indexed %d NX node names", len(nxNames))
	}

	// Walk scripts and mark deletions if we can't find a matching "<id>.img"
	// in NX (considering multiple zero-padding variants).
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
		id, perr := strconv.Atoi(base)
		if perr != nil {
			if verbose {
				log.Printf("Skipping non-numeric script file: %s", d.Name())
			}
			return nil
		}

		if !npcImgExists(nxNames, id) {
			toDelete = append(toDelete, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk scripts dir: %v", err)
	}

	if len(toDelete) == 0 {
		log.Println("No orphan NPC script files found.")
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

// npcImgExists checks if any plausible "<id>.img" filename exists in NX,
// accounting for leading zeros in NX data.
func npcImgExists(nxNames map[string]struct{}, id int) bool {
	candidates := []string{
		strconv.Itoa(id) + ".img",
		sprintfPad(id, 7) + ".img",
		sprintfPad(id, 8) + ".img",
		sprintfPad(id, 9) + ".img",
	}
	for _, c := range candidates {
		if _, ok := nxNames[strings.ToLower(c)]; ok {
			return true
		}
	}
	return false
}

func sprintfPad(id int, width int) string {
	// zero-pad to width
	s := strconv.Itoa(id)
	if len(s) >= width {
		return s
	}
	return strings.Repeat("0", width-len(s)) + s
}

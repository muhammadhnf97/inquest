package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/muhammadhnf97/inquest/internal"
)

var version = "1.0.0"

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")

	var path string
	flag.StringVar(&path, "path", ".", "Path to investigate")
	flag.StringVar(&path, "p", ".", "Path to investigate")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
	}

	if path != "" {
		var totalSize int64
		info := internal.ValidatePath(path)
		if info.IsDir {
			fmt.Printf("inquest v%s — investigating %s\n\n", version, path)

			suspects := internal.Gathering(path)
			guilty := internal.Investigate(suspects)

			if len(guilty) == 0 {
				fmt.Println("Inquest complete. All innocent.")
				return
			}
			fmt.Printf("%d of %d suspects stand accused.\n\n", len(guilty), len(suspects))

			for _, g := range guilty {
				fmt.Printf("%-9s %-13s %s  (evidence: %s)\n",
					internal.Humanize(g.Weight),
					g.Artifact,
					g.Path,
					filepath.Base(g.Evidence))
				totalSize += g.Weight
			}
			disk := internal.DiskFree(path)
			pct := float64(totalSize) / float64(disk.Total) * 100
			fmt.Println("")
			fmt.Printf("Reclaimable space: %s (%.1f%% of %s disk)\n", internal.Humanize(totalSize), pct, internal.Humanize(disk.Total))
			fmt.Println("Reclaim the space? (y/N)")

			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				fmt.Println("Reclaiming space...")
			}

		} else {
			fmt.Printf("Not a directory: %s\n", path)
		}
	} else {
		fmt.Println("No path provided")
	}
}

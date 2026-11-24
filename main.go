package main

import (
	"flag"
	"fmt"

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
			fmt.Println("INFO:     A directory, proceeding to gathering suspect...")
			fmt.Println("INFO:     Gathering charges...")
			suspects := internal.Gathering(path)

			fmt.Println("INFO:     Suspects found:", len(suspects))
			guilty := internal.Investigate(suspects)
			fmt.Println("INFO:     Guilty:", len(guilty))

			for _, g := range guilty {
				fmt.Println("Guilty:", g.Path)
				fmt.Println("Evidence:", g.Evidence)
				fmt.Println("Size:", internal.Humanize(g.Weight))
				fmt.Println()
				totalSize += g.Weight
			}
			fmt.Println("Reclaimable space:", internal.Humanize(totalSize))

		} else {
			fmt.Println("INFO:     Not a directory")
		}
	} else {
		fmt.Println("INFO:     No path provided")
	}
}

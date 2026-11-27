package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	"github.com/muhammadhnf97/inquest/internal"
)

var version = "0.1.0"

func main() {
	if info, ok := debug.ReadBuildInfo(); ok &&
		info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = strings.TrimPrefix(info.Main.Version, "v")
	}
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")

	var path string
	flag.StringVar(&path, "path", ".", "Path to investigate")
	flag.StringVar(&path, "p", ".", "Path to investigate")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run")
	flag.BoolVar(&dryRun, "n", false, "Dry run")

	var yes bool
	flag.BoolVar(&yes, "yes", false, "Execute without asking")
	flag.BoolVar(&yes, "y", false, "Execute without asking")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `inquest — your disk is under investigation

usage: inquest [--path <dir>] [--dry-run] [--yes] [--version]

flags:
  -p, --path <dir>   path to investigate (default ".")
  -n, --dry-run      judge only, execute nothing
  -y, --yes          execute all guilty without asking
  -v, --version      print version and exit
`)
	}
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var guilty []internal.Suspect

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments:", flag.Args())
		flag.Usage()
		os.Exit(2)
	}

	go func() {
		<-c
		fmt.Println("\nInquest interrupted")
		if len(guilty) > 0 {
			for _, g := range guilty {
				fmt.Println(g.Path)
			}
		}
		os.Exit(130)
	}()

	if showVersion {
		fmt.Println(version)
		return
	}

	if path != "" {
		var totalSize int64
		info := internal.ValidatePath(path)
		if info.IsDir {
			fmt.Printf("inquest v%s — investigating %s\n\n", version, path)

			suspects := internal.Gathering(path)
			guilty = internal.Investigate(suspects)

			if len(guilty) == 0 {
				fmt.Println("Inquest complete. All innocent.")
				return
			}
			fmt.Printf("%d of %d suspects stand accused.\n\n", len(guilty), len(suspects))

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SIZE\tTYPE\tAGE\tPATH")
			for _, g := range guilty {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					internal.Humanize(g.Weight),
					g.Artifact,
					internal.HumanizeAge(g.ModTime),
					g.Path)
				totalSize += g.Weight
			}
			w.Flush()
			disk := internal.DiskFree(path)
			pct := float64(totalSize) / float64(disk.Total) * 100
			fmt.Println("")
			fmt.Printf("Reclaimable space: %s (%.1f%% of %s disk)\n", internal.Humanize(totalSize), pct, internal.Humanize(disk.Total))
			if dryRun {
				fmt.Println("Dry run mode. No files will be deleted.")
				return
			}
			var response string
			if !yes {
				fmt.Println("Reclaim the space? (y/N)")
				fmt.Scanln(&response)
			}
			if yes || response == "y" || response == "Y" || response == "yes" || response == "YES" {
				fmt.Println("Reclaiming space...")

				var reclaimed int64
				for _, g := range guilty {
					if err := internal.Execute(g.Path); err != nil {
						fmt.Println("Failed to execute:", g.Path, "-", err)
					} else {
						fmt.Printf("Executed: %s (reclaimed %s)\n", g.Path, internal.Humanize(g.Weight))
						reclaimed += g.Weight
					}
				}
				if reclaimed == totalSize {
					fmt.Println("Sentence carried out. All executions complete.")
				} else {
					fmt.Printf("Sentence partially carried out. %s of %s reclaimed.\n",
						internal.Humanize(reclaimed), internal.Humanize(totalSize))
				}
			} else {
				fmt.Println("Guilty has been forgiven.")
			}

		} else {
			fmt.Printf("Not a directory: %s\n", path)
		}
	} else {
		fmt.Println("No path provided")
	}
}

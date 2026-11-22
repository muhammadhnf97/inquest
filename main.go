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
		info := internal.ValidatePath(path)
		if info.IsDir {
			fmt.Println("INFO     : A directory, proceeding to gathering suspect...")
			internal.Gathering(path)
		} else {
			fmt.Println("INFO     : Not a directory")
		}
	}
}

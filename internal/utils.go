package internal

import (
	"fmt"
	"os"
)

var charges map[string]string

type PathInfo struct {
	Path  string
	Size  int64
	IsDir bool
	Mode  os.FileMode
}

func ValidatePath(path string) PathInfo {
	info, err := os.Stat(path)
	if err != nil {
		return PathInfo{}
	}

	return PathInfo{
		Path:  path,
		Size:  info.Size(),
		IsDir: info.IsDir(),
		Mode:  info.Mode(),
	}
}

func Gathering(path string) {
	charges = map[string]string{
		"node_modules": "package.json",
		"venv":         "requirements.txt",
		".venv":        "requirements.txt",
	}
	fmt.Println("Gathering charges...")

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if info.IsDir() {
		fmt.Println("Is directory")
	}
}

func Investigate(path string) {

}

func Execute(path string) {

}

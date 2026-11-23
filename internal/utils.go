package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

type Suspect struct {
	Path     string // the accused dir — what gets deleted
	Artifact string // charged dir name (node_modules, .venv, ...)
	Evidence string // manifest that convicted it — set at trial, empty before
}

func Weight(path string, suspects []Suspect) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func Gathering(path string) []Suspect {
	charges = map[string]string{
		"node_modules": "package.json",
		"venv":         "requirements.txt",
		".venv":        "requirements.txt",
	}
	fmt.Println("INFO:     Gathering charges...")

	var suspects []Suspect
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == ".git" {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			return nil
		}

		if _, charged := charges[d.Name()]; charged {
			suspects = append(suspects, Suspect{
				Path:     path,
				Artifact: d.Name(),
			})
			return fs.SkipDir // outermost only — never descend into a suspect
		}
		return nil
	})

	return suspects
}

// Investigate tries each suspect: evidence present in the parent dir → guilty.
func Investigate(suspects []Suspect) []Suspect {
	var guilty []Suspect
	for _, s := range suspects {
		evidence := filepath.Join(filepath.Dir(s.Path), charges[s.Artifact])
		if _, err := os.Stat(evidence); err == nil {
			s.Evidence = evidence
			guilty = append(guilty, s)
		}
	}
	return guilty
}

func Execute(path string) {

}

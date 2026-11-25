package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func Gathering(path string) []Suspect {
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
			var modTime time.Time
			if info, err := d.Info(); err == nil {
				modTime = info.ModTime()
			}
			suspects = append(suspects, Suspect{
				Path:     path,
				Artifact: d.Name(),
				ModTime:  modTime,
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
			s.Weight = Weight(s.Path)
			guilty = append(guilty, s)
		}
	}
	return guilty
}

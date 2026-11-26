package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGathering(t *testing.T) {
	tests := []struct {
		name  string
		plant func(dir string) // builds the scene
		want  int              // expected suspect count
	}{
		{"artifact dir gathered", plantGuilty, 1},
		{"artifact without manifest still gathered", plantInnocent, 1},
		{"nested: outermost only", plantNested, 1},
		{"file named like artifact ignored", plantFileSuspect, 0},
		{"git dir never entered", plantGit, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			tc.plant(dir)
			if got := len(Gathering(dir)); got != tc.want {
				t.Errorf("got %d suspects, want %d", got, tc.want)
			}
		})
	}
}

func TestInvestigate(t *testing.T) {
	tests := []struct {
		name  string
		plant func(dir string)
		want  int // expected guilty count
	}{
		{"manifest present: convicted", plantGuilty, 1},
		{"no manifest: acquitted", plantInnocent, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			tc.plant(dir)
			suspects := Gathering(dir)
			guilty := Investigate(suspects)
			if len(guilty) != tc.want {
				t.Errorf("got %d guilty, want %d", len(guilty), tc.want)
			}
		})
	}
}

func TestInvestigateRecordsEvidence(t *testing.T) {
	dir := t.TempDir()
	plantGuilty(dir)

	guilty := Investigate(Gathering(dir))
	if len(guilty) != 1 {
		t.Fatalf("got %d guilty, want 1", len(guilty))
	}

	want := filepath.Join(dir, "package.json")
	if guilty[0].Evidence != want {
		t.Errorf("Evidence = %q, want %q", guilty[0].Evidence, want)
	}
}

func plantGuilty(dir string) {
	os.MkdirAll(filepath.Join(dir, "node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
}

func plantInnocent(dir string) {
	os.MkdirAll(filepath.Join(dir, "node_modules"), 0755) // no manifest
}

func plantNested(dir string) {
	os.MkdirAll(filepath.Join(dir, "node_modules/node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
}

func plantFileSuspect(dir string) {
	os.WriteFile(filepath.Join(dir, "node_modules"), []byte("not a dir"), 0644)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
}

func plantGit(dir string) {
	os.MkdirAll(filepath.Join(dir, ".git/node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, ".git/package.json"), []byte("{}"), 0644)
}

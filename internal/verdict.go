package internal

import "time"

var charges = map[string]string{
	"node_modules": "package.json",
	"venv":         "requirements.txt",
	".venv":        "requirements.txt",
	"target":       "Cargo.toml",
	".next":        "package.json",
	".nuxt":        "package.json",
	"vendor":       "go.mod",
}

type Suspect struct {
	Path     string    // the accused dir — what gets deleted
	Artifact string    // charged dir name (node_modules, .venv, ...)
	Evidence string    // manifest that convicted it — set at trial, empty before
	Weight   int64     // weight of the suspect
	ModTime  time.Time // last modified — dormancy signal
}

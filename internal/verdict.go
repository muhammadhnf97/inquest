package internal

var charges = map[string]string{
	"node_modules": "package.json",
	"venv":         "requirements.txt",
	".venv":        "requirements.txt",
}

type Suspect struct {
	Path     string // the accused dir — what gets deleted
	Artifact string // charged dir name (node_modules, .venv, ...)
	Evidence string // manifest that convicted it — set at trial, empty before
	Weight   int64  // weight of the suspect
}

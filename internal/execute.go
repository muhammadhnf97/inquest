package internal

import "os"

func Execute(path string) error {
	return os.RemoveAll(path)
}

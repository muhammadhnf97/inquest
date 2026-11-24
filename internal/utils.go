package internal

import "os"

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

package internal

import (
	"fmt"
	"os"
	"time"
)

func HumanizeAge(t time.Time) string {
	days := int(time.Since(t).Hours() / 24)
	switch {
	case days < 1:
		return "today"
	case days == 1:
		return "1 day ago"
	case days < 30:
		return fmt.Sprintf("%d days ago", days)
	case days < 60:
		return "1 month ago"
	case days < 365:
		return fmt.Sprintf("%d months ago", days/30)
	case days < 730:
		return "1 year ago"
	default:
		return fmt.Sprintf("%d years ago", days/365)
	}
}

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

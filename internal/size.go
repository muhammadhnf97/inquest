package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"
)

func Weight(path string) int64 {
	var total int64
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		total += info.Size()
		return nil

	})
	return total
}

func Humanize(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int(unit), 0

	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

type DiskInfo struct {
	Total int64
	Free  int64
	Used  int64
}

func DiskFree(path string) DiskInfo {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskInfo{-1, -1, -1}
	}
	total := int64(stat.Blocks) * int64(stat.Bsize)           // disk capacity
	free := int64(stat.Bavail) * int64(stat.Bsize)            // what you can use
	used := int64(stat.Blocks-stat.Bfree) * int64(stat.Bsize) // occupied
	return DiskInfo{total, free, used}
}

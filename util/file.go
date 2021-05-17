package util

import (
	"os"
	"path/filepath"
)

func CheckDir(path string) (string, bool) {
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", false
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || !fi.IsDir() {
		return "", false
	}

	return path, true
}

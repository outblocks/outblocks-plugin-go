package util

import (
	"os"
	"path/filepath"
)

func CheckDir(path string) (string, bool) {
	eval, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, false
	}

	fi, err := os.Stat(eval)
	if os.IsNotExist(err) || !fi.IsDir() {
		return path, false
	}

	return eval, true
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DirExists(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

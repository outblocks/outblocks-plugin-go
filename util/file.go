package util

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
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

func WriteFile(filename string, data []byte, perm fs.FileMode) error {
	err := os.WriteFile(filename, data, perm)
	if err != nil {
		return err
	}

	uidStr, ok1 := os.LookupEnv("SUDO_UID")
	gidStr, ok2 := os.LookupEnv("SUDO_GID")

	if ok1 && ok2 {
		uid, _ := strconv.Atoi(uidStr)
		gid, _ := strconv.Atoi(gidStr)

		return os.Chown(filename, uid, gid)
	}

	return nil
}

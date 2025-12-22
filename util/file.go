package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
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

func CheckMatch(name string, globs []glob.Glob) bool {
	for _, g := range globs {
		if g.Match(name) {
			return true
		}
	}

	return false
}

func WalkWithExclusions(dir string, excludes []string, fn func(path, rel string, info os.FileInfo) error) error {
	for i := range excludes {
		excludes[i] = filepath.FromSlash(excludes[i])
	}

	var (
		g                          glob.Glob
		excludeGlobs, includeGlobs []glob.Glob
		err                        error
	)

	for _, pat := range excludes {
		if pat != "" && pat[0] == '!' {
			g, err = glob.Compile(pat[1:])
			if err != nil {
				return fmt.Errorf("unable to parse inclusion '%s': %w", pat, err)
			}

			excludeGlobs = append(excludeGlobs, g)
		} else {
			g, err = glob.Compile(pat)
			if err != nil {
				return fmt.Errorf("unable to parse exclusion '%s': %w", pat, err)
			}

			includeGlobs = append(includeGlobs, g)
		}
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error encountered during file walk: %w", err)
		}

		relname, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("error relativizing file: %w", err)
		}

		isMatch := CheckMatch(relname, excludeGlobs)
		if isMatch && CheckMatch(relname, includeGlobs) {
			isMatch = false
		}

		if info.IsDir() {
			if isMatch {
				return filepath.SkipDir
			}

			return nil
		}

		if isMatch {
			return nil
		}

		return fn(path, relname, info)
	})
}

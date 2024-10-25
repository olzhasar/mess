package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// TODO: Make this configurable
var globs = []string{
	"*.pyc",
	"__pycache__",
	".mypy_cache",
	".pytest_cache",
	".ruff_cache",
	".tox",
	".nox",
	"node_modules",
}

func Clean(path string, verbose bool) (int, error) {
	counter := 0

	absPath, err := parsePath(path)
	if err != nil {
		return 0, err
	}

	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, glob := range globs {
			matched, err := filepath.Match(glob, d.Name())
			if err != nil {
				return err
			}
			if matched {
				counter++
				if verbose {
					fmt.Println("Removing", path)
				}
				os.RemoveAll(path)
				return fs.SkipDir
			}
		}

		return nil
	})

	return counter, err
}

package lib

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// TODO: Make this configurable
var patterns = []string{
	"*.pyc",
	"__pycache__",
	".mypy_cache",
	".pytest_cache",
	".ruff_cache",
	".tox",
	".nox",
	"node_modules",
}

type CleanOptions struct {
	Patterns  []string
	Verbose   bool
	Recursive bool
}

func getCleanPatterns(options CleanOptions) []string {
	if len(options.Patterns) > 0 {
		return options.Patterns
	}

	return patterns
}

func Clean(root string, options CleanOptions, stdout io.Writer) (int, error) {
	rootPath, err := parsePath(root)
	if err != nil {
		return 0, err
	}

	patterns := getCleanPatterns(options)

	counter := 0
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == rootPath {
			return nil
		}

		relPath, relErr := filepath.Rel(rootPath, path)
		if relErr != nil {
			panic("build relative path")
		}

		for _, pattern := range patterns {
			// First try the entry name only, if that doesn't work, try matching the relative path.
			// The latter is needed to cover patterns like /repo/node_modules/

			matched, err := filepath.Match(pattern, d.Name())
			if err != nil {
				return err
			}
			if !matched {
				matched, err = filepath.Match(pattern, relPath)
				if err != nil {
					return err
				}
			}

			if matched {
				if options.Verbose {
					fmt.Fprintln(stdout, "Removing ", relPath)
				}
				removeErr := os.RemoveAll(path)
				if removeErr != nil {
					return removeErr
				}
				counter++

				if d.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}

		if !options.Recursive && d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})

	return counter, err
}

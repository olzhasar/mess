package lib

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

type CleanOptions struct {
	Patterns  []string
	Verbose   bool
	Recursive bool
	CalcFreed bool
}

type CleanResult struct {
	BytesFreed uint64
	Count      uint64
}

func getCleanPatterns(options CleanOptions) ([]string, error) {
	if len(options.Patterns) > 0 {
		return options.Patterns, nil
	}

	patterns, err := LoadPatterns()
	if err != nil {
		return nil, err
	}

	return patterns, nil
}

func Clean(root string, options CleanOptions, stdout io.Writer) (CleanResult, error) {
	result := CleanResult{}

	rootPath, err := parsePath(root)
	if err != nil {
		return result, err
	}

	patterns, err := getCleanPatterns(options)
	if err != nil {
		return result, err
	}

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

				var size uint64 = 0
				if options.CalcFreed {
					size, err = diskUsage(path)
					if err != nil {
						return err
					}
				}

				removeErr := os.RemoveAll(path)
				if removeErr != nil {
					return removeErr
				}
				result.Count++
				result.BytesFreed += size

				if d.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}

		if !options.Recursive && d.IsDir() {
			return fs.SkipDir
		}

		return nil
	})

	return result, err
}

func diskUsage(root string) (uint64, error) {
	var total uint64 = 0

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if ok {
			total += uint64(stat.Blocks) * 512
		}

		return nil
	})

	return total, err
}

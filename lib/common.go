package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func parsePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Cannot open %s, make sure the path exists", path)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("Cannot open %s, make sure the path exists", absPath)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", absPath)
	}

	return absPath, nil
}

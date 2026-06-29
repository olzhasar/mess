package lib

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed PATTERNS
var builtInPatternsFile string

func findPatternsFile() string {
	candidates := make([]string, 0)

	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(homeDir, ".mess_patterns"))
	}

	if configDir, err := os.UserConfigDir(); err == nil {
		candidates = append(candidates, filepath.Join(configDir, "mess", "patterns"))
	}

	for _, path := range candidates {

		stat, err := os.Stat(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "error accessing %s:\n%v\n", path, err)
			}
			continue
		}
		if stat.Mode().IsDir() {
			fmt.Fprintf(os.Stderr, "found %s, but it's a directory\n", path)
			continue
		}

		return path
	}

	return ""
}

func readPatterns(reader io.Reader) ([]string, error) {
	patterns := make([]string, 0)
	existing := make(map[string]bool, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		_, ok := existing[line]
		if ok {
			continue
		}
		existing[line] = true

		fmt.Println(line)

		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(patterns) == 0 {
		return nil, errors.New("empty file")
	}

	return patterns, nil
}

func readPatternsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readPatterns(file)
}

func loadPatterns() ([]string, error) {
	filePath := findPatternsFile()
	if filePath == "" {
		return readPatterns(strings.NewReader(builtInPatternsFile))
	}

	fmt.Fprintln(os.Stderr, "Using patterns from file:", filePath)

	return readPatternsFromFile(filePath)
}

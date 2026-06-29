package lib

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBuiltInPatternsAreValid(t *testing.T) {
	t.Parallel()

	patterns, err := readPatterns(strings.NewReader(builtInPatternsFile))
	if err != nil {
		t.Fatalf("readPatterns: %v", err)
	}

	for _, pattern := range patterns {
		if _, err := filepath.Match(pattern, "example"); err != nil {
			t.Fatalf("invalid pattern %q: %v", pattern, err)
		}
	}
}

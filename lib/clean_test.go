package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanRemovesGlobs(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	paths := []string{
		filepath.Join(root, "node_modules", "foo.js"),
		filepath.Join(root, "__pycache__", "bar.pyc"),
		filepath.Join(root, ".mypy_cache", "baz"),
		filepath.Join(root, "keep", "keep.txt"),
		filepath.Join(root, "test.pyc"),
	}

	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	removed, err := Clean(root, false)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}

	if removed != 4 {
		t.Fatalf("expected 4 items removed, got %d", removed)
	}

	removedPaths := []string{
		filepath.Join(root, "node_modules"),
		filepath.Join(root, "__pycache__"),
		filepath.Join(root, ".mypy_cache"),
		filepath.Join(root, "test.pyc"),
	}
	for _, path := range removedPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed", path)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "keep", "keep.txt")); err != nil {
		t.Fatalf("expected keep file to remain: %v", err)
	}
}

func TestCleanInvalidPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := Clean(file, false)
	if err == nil {
		t.Fatalf("expected error for file path")
	}
}

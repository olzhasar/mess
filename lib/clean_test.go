package lib_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/olzhasar/mess/lib"
)

func TestMain(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	mkDirs(t, root,
		[]string{
			"node_modules",
			"foo",
			"foo/bar",
			"foo/bar/node_modules",
		})
	mkFiles(t, root,
		[]string{
			"main.pyc",
			"README.md",
			"foo/specific",
			"foo/bar/test.pyc",
		})

	got, err := lib.Clean(root, lib.CleanOptions{Recursive: true, Patterns: []string{
		"*.pyc",
		"node_modules",
		"foo/specific",
		"README", // partial match should be excluded
	}}, nil)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	assertPathDeleted(t, root, "main.pyc")
	assertPathDeleted(t, root, "foo/specific")
	assertPathDeleted(t, root, "foo/bar/test.pyc")
	assertPathDeleted(t, root, "node_modules")
	assertPathDeleted(t, root, "foo/bar/node_modules")
	assertPathExists(t, root, "README.md")
	assertPathExists(t, root, "foo")
	assertPathExists(t, root, "foo/bar")

	assertDeletedCount(t, 5, int(got.Count))
}

func TestNoRecurse(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mkDirs(t, root, []string{
		"foo",
		"foo/bar",
		"bar", // should be deleted
	})
	mkFiles(t, root, []string{
		"main.pyc",
		"foo/test.pyc", // should stay
	})

	got, err := lib.Clean(root, lib.CleanOptions{
		Recursive: false,
		Patterns: []string{
			"*.pyc",
			"bar",
		}}, nil)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}
	assertPathDeleted(t, root, "main.pyc")
	assertPathDeleted(t, root, "bar")
	assertPathExists(t, root, "foo/bar")
	assertPathExists(t, root, "foo/test.pyc")

	assertDeletedCount(t, 2, int(got.Count))
}

func TestCalculatesFreed(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mkDirs(t, root, []string{
		"node_modules",
		"node_modules/package",
	})

	contents := bytes.Repeat([]byte("b"), 1<<20)
	mkFileWithContent(t, root, "node_modules/package/index.js", contents)

	got, err := lib.Clean(root, lib.CleanOptions{
		Recursive: true,
		CalcFreed: true,
		Patterns:  []string{"node_modules"},
	}, nil)
	if err != nil {
		t.Fatalf("Clean: %v", err)
	}

	assertPathDeleted(t, root, "node_modules")
	assertDeletedCount(t, 1, int(got.Count))
	if got.BytesFreed < uint64(len(contents)) {
		t.Fatalf("assert failed: want at least %d bytes freed, got %d", len(contents), got.BytesFreed)
	}
}

func TestInvalidPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := lib.Clean(file, lib.CleanOptions{Verbose: false}, nil)
	if err == nil {
		t.Fatalf("expected error for file path")
	}
}

func mkDirs(t *testing.T, root string, paths []string) {
	t.Helper()

	for _, path := range paths {
		fullPath := filepath.Join(root, path)
		if err := os.Mkdir(fullPath, 0o755); err != nil {
			t.Fatal("Failed to create directory", fullPath, err)
		}
	}
}

func mkFiles(t *testing.T, root string, paths []string) {
	t.Helper()

	for _, path := range paths {
		fullPath := filepath.Join(root, path)
		if err := os.WriteFile(fullPath, []byte{}, 0o644); err != nil {
			t.Fatal("Failed to created a file", fullPath, err)
		}
	}
}

func mkFileWithContent(t *testing.T, root string, path string, contents []byte) {
	t.Helper()

	fullPath := filepath.Join(root, path)
	if err := os.WriteFile(fullPath, contents, 0o644); err != nil {
		t.Fatal("Failed to create a file", fullPath, err)
	}
}

func assertPathDeleted(tb testing.TB, root string, path string) {
	tb.Helper()

	fullPath := filepath.Join(root, path)

	_, err := os.Stat(fullPath)
	if err == nil {
		tb.Fatalf("assert failed: path %s exists", fullPath)
	}
}

func assertPathExists(tb testing.TB, root string, path string) {
	tb.Helper()

	fullPath := filepath.Join(root, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		tb.Fatalf("assert failed: path %s does not exist", fullPath)
	}
}

func assertDeletedCount(tb testing.TB, want int, got int) {
	tb.Helper()

	if want != got {
		tb.Fatalf("assert failed: want %d removed paths, got %d", want, got)
	}
}

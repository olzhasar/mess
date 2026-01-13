package lib

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func TestFindGitReposBasic(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	repo1 := filepath.Join(root, "repo1")
	repo2 := filepath.Join(root, "repo2")
	if err := os.MkdirAll(repo1, 0o755); err != nil {
		t.Fatalf("mkdir repo1: %v", err)
	}
	if err := os.MkdirAll(repo2, 0o755); err != nil {
		t.Fatalf("mkdir repo2: %v", err)
	}

	initGitRepo(t, repo1)
	initGitRepo(t, repo2)

	results, err := FindGitRepos(root, false, 0, 0)
	if err != nil {
		t.Fatalf("FindGitRepos: %v", err)
	}

	sort.Strings(results)
	expected := []string{repo1, repo2}
	sort.Strings(expected)

	if len(results) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(results))
	}
	for i, result := range results {
		if result != expected[i] {
			t.Fatalf("expected %s, got %s", expected[i], result)
		}
	}
}

func TestFindGitReposDirty(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repoDirty := filepath.Join(root, "repo-dirty")
	repoClean := filepath.Join(root, "repo-clean")

	if err := os.MkdirAll(repoDirty, 0o755); err != nil {
		t.Fatalf("mkdir repoDirty: %v", err)
	}
	if err := os.MkdirAll(repoClean, 0o755); err != nil {
		t.Fatalf("mkdir repoClean: %v", err)
	}

	initGitRepo(t, repoDirty)
	initGitRepo(t, repoClean)

	commitFile(t, repoDirty, "file.txt", "clean", time.Now())
	commitFile(t, repoClean, "file.txt", "clean", time.Now())

	if err := os.WriteFile(filepath.Join(repoDirty, "file.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	results, err := FindGitRepos(root, true, 0, 0)
	if err != nil {
		t.Fatalf("FindGitRepos: %v", err)
	}

	if len(results) != 1 || results[0] != repoDirty {
		t.Fatalf("expected only dirty repo %s, got %v", repoDirty, results)
	}
}

func TestFindGitReposOlder(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repoOld := filepath.Join(root, "repo-old")
	repoNew := filepath.Join(root, "repo-new")

	if err := os.MkdirAll(repoOld, 0o755); err != nil {
		t.Fatalf("mkdir repoOld: %v", err)
	}
	if err := os.MkdirAll(repoNew, 0o755); err != nil {
		t.Fatalf("mkdir repoNew: %v", err)
	}

	initGitRepo(t, repoOld)
	initGitRepo(t, repoNew)

	commitFile(t, repoOld, "file.txt", "old", time.Now().AddDate(0, 0, -10))
	commitFile(t, repoNew, "file.txt", "new", time.Now())

	results, err := FindGitRepos(root, false, 5, 0)
	if err != nil {
		t.Fatalf("FindGitRepos: %v", err)
	}

	if len(results) != 1 || results[0] != repoOld {
		t.Fatalf("expected only old repo %s, got %v", repoOld, results)
	}
}

func TestFindGitReposInvalidPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := FindGitRepos(file, false, 0, 0)
	if err == nil {
		t.Fatalf("expected error for file path")
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func commitFile(t *testing.T, dir, name, contents string, commitTime time.Time) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	runGit(t, dir, "add", name)

	cmd := exec.Command("git", "-C", dir, "commit", "-m", "commit")
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+commitTime.Format(time.RFC3339),
		"GIT_COMMITTER_DATE="+commitTime.Format(time.RFC3339),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git commit: %v: %s", err, output)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

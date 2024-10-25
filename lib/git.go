package lib

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Check func(path string, d fs.DirEntry) bool

func traverse(path string, preChecks []Check, postChecks []Check, limitConcurrency int) ([]string, error) {
	results := make([]string, 0)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println("Failed to access path", path)
			return fs.SkipDir
		}

		for _, check := range preChecks {
			if !check(path, d) {
				return nil
			}
		}

		results = append(results, path)
		return fs.SkipDir
	})

	if err != nil {
		return nil, err
	}

	if postChecks == nil || len(postChecks) == 0 {
		return results, nil
	}

	wg := sync.WaitGroup{}
	ch := make(chan string, len(results))

	if limitConcurrency < 1 {
		limitConcurrency = 100
	}
	sem := make(chan struct{}, limitConcurrency)

	for _, result := range results {
		wg.Add(1)
		go func(result string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			for _, check := range postChecks {
				if !check(result, nil) {
					return
				}
			}
			ch <- result
		}(result)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results = make([]string, 0)
	for result := range ch {
		results = append(results, result)
	}

	return results, nil
}

func isGitRepo(path string, d fs.DirEntry) bool {
	if !d.IsDir() {
		return false
	}

	_, err := os.Stat(filepath.Join(path, ".git"))
	if err == nil {
		return true
	}

	return false
}

func hasUncommittedChanges(path string, _ fs.DirEntry) bool {
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Failed to check git status in", path)
		return false
	}
	if len(output) > 0 {
		return true
	}
	return false
}

func lastCommitOlderThan(days int) Check {
	earliestTime := time.Now().AddDate(0, 0, -days)
	dateISO := earliestTime.Format("2006-01-02")

	return func(path string, _ fs.DirEntry) bool {
		cmd := exec.Command("git", "-C", path, "log", "-1", "--since="+dateISO)
		output, err := cmd.Output()
		if err != nil {
			// usually means there are no commits
			return false
		}

		if len(output) == 0 {
			return true
		}

		return false
	}
}

func FindGitRepos(path string, dirty bool, older int, limitConcurrency int) ([]string, error) {
	path, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	preChecks := []Check{isGitRepo}
	postChecks := []Check{}

	if dirty {
		postChecks = append(postChecks, hasUncommittedChanges)
	}

	if older > 0 {
		postChecks = append(postChecks, lastCommitOlderThan(older))
	}

	results, err := traverse(path, preChecks, postChecks, limitConcurrency)
	if err != nil {
		return nil, err
	}

	return results, nil
}

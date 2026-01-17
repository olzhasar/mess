package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	_, dirty, err := getDirtyStatus(path)
	if err != nil {
		log.Println("Failed to check git status in", path)
		return false
	}
	return dirty
}

type DirtyStats struct {
	ChangedFiles   int
	AddedLines     int
	DeletedLines   int
	UntrackedFiles int
}

type DirtyRepoStats struct {
	Path  string
	Stats DirtyStats
	Err   error
}

func GetDirtyStats(path string) (DirtyStats, error) {
	stats, dirty, err := getDirtyStatus(path)
	if err != nil {
		return stats, err
	}
	if !dirty {
		return stats, nil
	}

	added, deleted, err := diffNumstat(path, false)
	if err != nil {
		return stats, err
	}
	addedCached, deletedCached, err := diffNumstat(path, true)
	if err != nil {
		return stats, err
	}

	stats.AddedLines = added + addedCached
	stats.DeletedLines = deleted + deletedCached

	return stats, nil
}

func FindDirtyGitReposWithStats(path string, older int, limitConcurrency int, includeStats bool) ([]DirtyRepoStats, error) {
	results, err := FindGitRepos(path, false, older, limitConcurrency)
	if err != nil {
		return nil, err
	}

	if limitConcurrency < 1 {
		limitConcurrency = 100
	}
	sem := make(chan struct{}, limitConcurrency)
	wg := sync.WaitGroup{}
	repos := make([]DirtyRepoStats, 0, len(results))
	mu := sync.Mutex{}

	for _, result := range results {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			stats, dirty, err := getDirtyStatus(repo)
			if err != nil {
				log.Println("Failed to check git status in", repo)
				return
			}
			if !dirty {
				return
			}

			repoStats := DirtyRepoStats{
				Path:  repo,
				Stats: stats,
			}
			if includeStats {
				added, deleted, err := diffNumstat(repo, false)
				if err != nil {
					repoStats.Err = err
				} else {
					addedCached, deletedCached, err := diffNumstat(repo, true)
					if err != nil {
						repoStats.Err = err
					} else {
						repoStats.Stats.AddedLines = added + addedCached
						repoStats.Stats.DeletedLines = deleted + deletedCached
					}
				}
			}

			mu.Lock()
			repos = append(repos, repoStats)
			mu.Unlock()
		}(result)
	}

	wg.Wait()
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Path < repos[j].Path
	})
	return repos, nil
}

func getDirtyStatus(path string) (DirtyStats, bool, error) {
	var stats DirtyStats

	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return stats, false, fmt.Errorf("git status in %s: %w", path, err)
	}

	output = bytes.TrimSpace(output)
	if len(output) == 0 {
		return stats, false, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		stats.ChangedFiles++
		if strings.HasPrefix(line, "??") {
			stats.UntrackedFiles++
		}
	}
	if err := scanner.Err(); err != nil {
		return stats, false, fmt.Errorf("parse status in %s: %w", path, err)
	}

	return stats, true, nil
}

func diffNumstat(path string, cached bool) (int, int, error) {
	args := []string{"-C", path, "diff", "--numstat"}
	if cached {
		args = append(args, "--cached")
	}
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("git diff in %s: %w", path, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	added := 0
	deleted := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		added += parseNumstat(fields[0])
		deleted += parseNumstat(fields[1])
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("parse numstat in %s: %w", path, err)
	}

	return added, deleted, nil
}

func parseNumstat(value string) int {
	if value == "-" {
		return 0
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return number
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

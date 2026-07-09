package scanner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitstat/internal/model"
)

func gitExec(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %v in %s: %w", args, repoPath, err)
	}
	return strings.TrimRight(string(out), "\n\r "), nil
}

func ScanCommitsByRange(repoPath string, since, until time.Time) ([]model.Commit, error) {
	args := []string{"log", "--format=---GITSTAT_COMMIT---%n%H%n%an%n%ae%n%ci%n%s", "--numstat"}
	if !since.IsZero() {
		args = append(args, "--since="+since.Format("2006-01-02 15:04:05 -0700"))
	}
	if !until.IsZero() {
		args = append(args, "--until="+until.Format("2006-01-02 15:04:05 -0700"))
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("git log stdout pipe in %s: %w", repoPath, err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("git log start in %s: %w", repoPath, err)
	}

	commits, parseErr := parseGitLog(stdout)
	waitErr := cmd.Wait()

	if parseErr != nil {
		return nil, parseErr
	}
	return commits, waitErr
}

func parseGitLog(r io.Reader) ([]model.Commit, error) {
	scanner := bufio.NewScanner(r)
	const marker = "---GITSTAT_COMMIT---"

	// skip to first marker
	for scanner.Scan() {
		if scanner.Text() == marker {
			break
		}
	}

	var commits []model.Commit

	for {
		if !scanner.Scan() {
			break
		}
		hash := scanner.Text()

		if !scanner.Scan() {
			break
		}
		author := scanner.Text()

		if !scanner.Scan() {
			break
		}
		email := scanner.Text()

		if !scanner.Scan() {
			break
		}
		dateStr := scanner.Text()

		if !scanner.Scan() {
			break
		}
		subject := scanner.Text()

		commitTime, _ := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if commitTime.IsZero() {
			t, _ := time.Parse("2006-01-02 15:04:05", dateStr)
			commitTime = t.In(time.Local)
		}

		var additions, deletions int
		var files []model.FileStat
		for scanner.Scan() {
			line := scanner.Text()
			if line == marker {
				break
			}
			if line == "" {
				continue
			}
			parts := strings.Split(line, "\t")
			if len(parts) >= 3 {
				add, errA := strconv.Atoi(parts[0])
				del, errD := strconv.Atoi(parts[1])
				fpath := parts[2]
				if errA == nil {
					additions += add
				}
				if errD == nil {
					deletions += del
				}
				if fpath != "" && (errA == nil || errD == nil) {
					files = append(files, model.FileStat{
						Path:      fpath,
						Additions: add,
						Deletions: del,
					})
				}
			}
		}
		if files == nil {
			files = []model.FileStat{}
		}

		commits = append(commits, model.Commit{
			Hash:      hash,
			Author:    author,
			Email:     email,
			Date:      commitTime,
			Message:   subject,
			Additions: additions,
			Deletions: deletions,
			Files:     files,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return commits, nil
}


func GetHEADHash(repoPath string) (string, error) {
	return gitExec(repoPath, "rev-parse", "HEAD")
}

func IsAncestor(repoPath, hash string) bool {
	_, err := gitExec(repoPath, "merge-base", "--is-ancestor", hash, "HEAD")
	return err == nil
}

func ScanCommitsSince(repoPath, sinceHash string) ([]model.Commit, error) {
	args := []string{"log", sinceHash + "..HEAD", "--format=---GITSTAT_COMMIT---%n%H%n%an%n%ae%n%ci%n%s", "--numstat"}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("git log stdout pipe in %s: %w", repoPath, err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("git log start in %s: %w", repoPath, err)
	}
	commits, parseErr := parseGitLog(stdout)
	waitErr := cmd.Wait()
	if parseErr != nil {
		return nil, parseErr
	}
	return commits, waitErr
}

func DiscoverRepos(rootPath string) ([]model.Repository, error) {
	var repos []model.Repository

	gitPath := filepath.Join(rootPath, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		meta, err := ScanMetadata(rootPath)
		if err == nil {
			repos = append(repos, meta)
		}
	}

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		repoPath := filepath.Join(rootPath, entry.Name())
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
			continue
		}
		meta, err := ScanMetadata(repoPath)
		if err != nil {
			continue
		}
		repos = append(repos, meta)
	}
	return repos, nil
}

func ScanMetadata(repoPath string) (model.Repository, error) {
	currentBranch, err := gitExec(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Printf("warning: failed to get current branch for %s: %v", repoPath, err)
	}
	userEmail, err := gitExec(repoPath, "config", "user.email")
	if err != nil {
		log.Printf("warning: failed to get user email for %s: %v", repoPath, err)
	}

	var lastCommitTime string
	out, err := gitExec(repoPath, "log", "-1", "--format=%ci")
	if err == nil {
		if t, e := time.Parse("2006-01-02 15:04:05 -0700", out); e == nil {
			lastCommitTime = t.Format("2006-01-02 15:04:05")
		}
	} else {
		log.Printf("warning: failed to get last commit time for %s: %v", repoPath, err)
	}

	return model.Repository{
		Path:           repoPath,
		Name:           filepath.Base(repoPath),
		CurrentBranch:  currentBranch,
		UserEmail:      userEmail,
		LastCommitTime: lastCommitTime,
	}, nil
}

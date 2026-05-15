package scanner

import (
	"os"
	"path/filepath"
	"time"

	"gitstat/internal/model"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func ScanDirectory(path string, timeRange string) ([]model.Repository, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var repos []model.Repository

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		repoPath := filepath.Join(path, entry.Name())
		gitPath := filepath.Join(repoPath, ".git")

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			continue
		}

		repo, err := scanSingleRepo(repoPath, timeRange)
		if err != nil {
			continue // 跳过错误的仓库
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

func scanSingleRepo(path string, timeRange string) (model.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return model.Repository{}, err
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return model.Repository{}, err
	}

	var commits []model.Commit
	cutoffTime := calculateCutoffTime(timeRange)

	err = iter.ForEach(func(c *object.Commit) error {
		if c.Committer.When.Before(cutoffTime) {
			return nil
		}

		stats, _ := c.Stats()
		var additions, deletions int
		for _, stat := range stats {
			additions += stat.Addition
			deletions += stat.Deletion
		}

		commit := model.Commit{
			Hash:      c.Hash.String(),
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Date:      c.Committer.When,
			Message:   c.Message,
			Additions: additions,
			Deletions: deletions,
		}

		commits = append(commits, commit)
		return nil
	})

	return model.Repository{
		Path:    path,
		Name:    filepath.Base(path),
		Commits: commits,
	}, nil
}

func calculateCutoffTime(timeRange string) time.Time {
	switch timeRange {
	case "7d":
		return time.Now().AddDate(0, 0, -7)
	case "30d":
		return time.Now().AddDate(0, 0, -30)
	case "90d":
		return time.Now().AddDate(0, 0, -90)
	default:
		return time.Time{} // 全部历史
	}
}

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
	// 计算时间范围
	var startDate, endDate time.Time
	now := time.Now()

	switch timeRange {
	case "1d":
		startDate = now.AddDate(0, 0, -1)
		endDate = now
	case "7d":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "30d":
		startDate = now.AddDate(0, 0, -30)
		endDate = now
	case "90d":
		startDate = now.AddDate(0, 0, -90)
		endDate = now
	default: // all 或其他
		startDate = time.Time{} // 零值表示不限制
		endDate = now
	}

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

		repo, err := scanSingleRepo(repoPath, startDate, endDate)
		if err != nil {
			continue // 跳过错误的仓库
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

func scanSingleRepo(path string, startDate time.Time, endDate time.Time) (model.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return model.Repository{}, err
	}

	// 获取当前分支
	var currentBranch string
	head, err := repo.Head()
	if err == nil {
		currentBranch = head.Name().Short()
	}

	// 获取 Git 用户邮箱（项目级优先，否则全局）
	var userEmail string
	config, err := repo.Config()
	if err == nil && config.User.Email != "" {
		userEmail = config.User.Email
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return model.Repository{}, err
	}

	var commits []model.Commit
	var lastCommitTime time.Time

	err = iter.ForEach(func(c *object.Commit) error {
		commitTime := c.Committer.When

		// 如果提交时间晚于 endDate，跳过
		if !endDate.IsZero() && commitTime.After(endDate) {
			return nil
		}

		// 如果提交时间早于 startDate，停止扫描
		if !startDate.IsZero() && commitTime.Before(startDate) {
			return nil
		}

		// 记录最后提交时间
		if lastCommitTime.IsZero() || commitTime.After(lastCommitTime) {
			lastCommitTime = commitTime
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
			Date:      commitTime,
			Message:   c.Message,
			Additions: additions,
			Deletions: deletions,
		}

		commits = append(commits, commit)
		return nil
	})

	return model.Repository{
		Path:           path,
		Name:           filepath.Base(path),
		CurrentBranch:  currentBranch,
		LastCommitTime: lastCommitTime.Format("2006-01-02 15:04:05"),
		UserEmail:      userEmail,
		Commits:        commits,
	}, nil
}

func calculateCutoffTime(timeRange string) time.Time {
	switch timeRange {
	case "1d":
		return time.Now().AddDate(0, 0, -1)
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

// DiscoverRepos 发现目录下的所有 Git 仓库（仅元数据，不扫描提交）
func DiscoverRepos(rootPath string) ([]model.Repository, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var repos []model.Repository
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		repoPath := filepath.Join(rootPath, entry.Name())
		gitPath := filepath.Join(repoPath, ".git")

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
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

// ScanMetadata 扫描单个仓库元数据（不含提交）
func ScanMetadata(repoPath string) (model.Repository, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return model.Repository{}, err
	}

	var currentBranch string
	head, err := repo.Head()
	if err == nil {
		currentBranch = head.Name().Short()
	}

	var userEmail string
	config, err := repo.Config()
	if err == nil && config.User.Email != "" {
		userEmail = config.User.Email
	}

	// 获取最后提交时间
	var lastCommitTime string
	head, err = repo.Head()
	if err == nil {
		commit, err := repo.CommitObject(head.Hash())
		if err == nil {
			lastCommitTime = commit.Committer.When.Format("2006-01-02 15:04:05")
		}
	}

	return model.Repository{
		Path:           repoPath,
		Name:           filepath.Base(repoPath),
		CurrentBranch:  currentBranch,
		UserEmail:      userEmail,
		LastCommitTime: lastCommitTime,
	}, nil
}

// ScanCommitsByRange 扫描指定时间范围的提交
func ScanCommitsByRange(repoPath string, startDate, endDate time.Time) ([]model.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var commits []model.Commit
	err = iter.ForEach(func(c *object.Commit) error {
		commitTime := c.Committer.When

		if !endDate.IsZero() && commitTime.After(endDate) {
			return nil
		}

		if !startDate.IsZero() && commitTime.Before(startDate) {
			return nil
		}

		stats, _ := c.Stats()
		var additions, deletions int
		for _, stat := range stats {
			additions += stat.Addition
			deletions += stat.Deletion
		}

		commits = append(commits, model.Commit{
			Hash:      c.Hash.String(),
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Date:      commitTime,
			Message:   c.Message,
			Additions: additions,
			Deletions: deletions,
		})
		return nil
	})

	return commits, nil
}

// ScanIncremental 增量扫描指定时间范围的提交
// startDate: 扫描起始时间（更早的时间）
// endDate: 扫描结束时间（更晚的时间，通常是缓存的 EarliestDate）
func ScanIncremental(path string, startDate time.Time, endDate time.Time) ([]model.Commit, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var commits []model.Commit

	err = iter.ForEach(func(c *object.Commit) error {
		commitTime := c.Committer.When

		// 如果提交时间晚于 endDate，跳过（还未到达扫描范围）
		if !endDate.IsZero() && commitTime.After(endDate) {
			return nil
		}

		// 如果提交时间早于 startDate，停止扫描（已超出范围）
		if !startDate.IsZero() && commitTime.Before(startDate) {
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
			Date:      commitTime,
			Message:   c.Message,
			Additions: additions,
			Deletions: deletions,
		}

		commits = append(commits, commit)
		return nil
	})

	return commits, nil
}

package model

import "time"

type Repository struct {
	Path           string   `json:"path"`
	Name           string   `json:"name"`
	CurrentBranch  string   `json:"currentBranch"`
	LastCommitTime string   `json:"lastCommitTime"`
	UserEmail      string   `json:"userEmail"`
	Commits        []Commit `json:"commits"`
}

type Commit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
	Files     []string  `json:"files"`
}

type ScanRequest struct {
	Path      string `json:"path"`
	TimeRange string `json:"timeRange"` // e.g., "7d", "30d", "90d", "all"
}

type ScanResponse struct {
	Repositories []Repository `json:"repositories"`
	TotalCommits int          `json:"totalCommits"`
}

type OverviewStats struct {
	TotalCommits    int `json:"totalCommits"`
	TotalAdditions  int `json:"totalAdditions"`
	TotalDeletions  int `json:"totalDeletions"`
	ActiveAuthors   int `json:"activeAuthors"`
	RepositoryCount int `json:"repositoryCount"`
}

type AuthorDailyStats struct {
	Author    string          `json:"author"`
	Email     string          `json:"email"`
	Commits   int             `json:"commits"`
	Additions int             `json:"additions"`
	Deletions int             `json:"deletions"`
	IsMe      bool            `json:"isMe"`
	DailyData []DayCommitData `json:"dailyData,omitempty"`
}

type DayCommitData struct {
	Date      string `json:"date"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type RepositoryDailyStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string             `json:"currentBranch"`
	LastCommitTime string             `json:"lastCommitTime"`
	Authors        []AuthorDailyStats `json:"authors"`
}

// 周期聚合数据结构 (周/月/年)
type PeriodCommitData struct {
	Period    string `json:"date"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type AuthorPeriodStats struct {
	Author     string            `json:"author"`
	Email      string            `json:"email"`
	Commits    int               `json:"commits"`
	Additions  int               `json:"additions"`
	Deletions  int               `json:"deletions"`
	IsMe       bool              `json:"isMe"`
	PeriodData []PeriodCommitData `json:"dailyData,omitempty"`
}

type RepositoryPeriodStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string             `json:"currentBranch"`
	LastCommitTime string             `json:"lastCommitTime"`
	Authors        []AuthorPeriodStats `json:"authors"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// 开发者排行榜项
type AuthorRankItem struct {
	Author        string  `json:"author"`
	Email         string  `json:"email"`
	Commits       int     `json:"commits"`
	Additions     int     `json:"additions"`
	Deletions     int     `json:"deletions"`
	NetChange     int     `json:"netChange"` // 新增-删除
	IsMe          bool    `json:"isMe"`
	AvgCommitSize float64 `json:"avgCommitSize"` // 平均提交大小
}

// 活动热力图数据点
type ActivityHeatmapPoint struct {
	DayOfWeek   int `json:"dayOfWeek"` // 0=周日, 1=周一...6=周六
	Hour        int `json:"hour"`      // 0-23
	CommitCount int `json:"commitCount"`
}

// 仓库对比数据
type RepoComparison struct {
	RepoName         string  `json:"repoName"`
	RepoPath         string  `json:"repoPath"`
	Commits          int     `json:"commits"`
	Authors          int     `json:"authors"`
	Additions        int     `json:"additions"`
	Deletions        int     `json:"deletions"`
	LastCommitTime   string  `json:"lastCommitTime"`
	ActiveDays       int     `json:"activeDays"`       // 活跃天数
	AvgCommitsPerDay float64 `json:"avgCommitsPerDay"` // 日均提交数
}

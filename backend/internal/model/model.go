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

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

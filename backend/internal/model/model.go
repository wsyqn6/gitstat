package model

import "time"

type Repository struct {
	Path    string   `json:"path"`
	Name    string   `json:"name"`
	Commits []Commit `json:"commits"`
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

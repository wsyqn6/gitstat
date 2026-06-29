package store

import (
	"testing"
	"time"

	"gitstat/internal/model"
)

func TestStoreOperations(t *testing.T) {
	s := &Store{Repos: make(map[string]*RepoCache)}

	repos := []model.Repository{
		{
			Name: "test-repo",
			Path: "/path/to/test-repo",
			Commits: []model.Commit{
				{
					Hash:      "abc123",
					Author:    "Test User",
					Email:     "test@example.com",
					Date:      time.Now(),
					Message:   "Test commit",
					Additions: 10,
					Deletions: 5,
				},
			},
		},
	}

	s.RegisterRepos(repos)
	s.SetRepoCommits(repos[0].Path, repos[0].Commits)

	if len(s.GetRepositories()) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(s.GetRepositories()))
	}

	repos = s.GetRepositories()
	if len(repos[0].Commits) != 1 {
		t.Errorf("Expected 1 commit, got %d", len(repos[0].Commits))
	}

	if repos[0].Commits[0].Author != "Test User" {
		t.Errorf("Expected author 'Test User', got '%s'", repos[0].Commits[0].Author)
	}
}

func TestStoreEmpty(t *testing.T) {
	s := &Store{Repos: make(map[string]*RepoCache)}

	if len(s.GetRepositories()) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(s.GetRepositories()))
	}

	repos2 := s.GetRepositories()
	for _, r := range repos2 {
		if len(r.Commits) != 0 {
			t.Errorf("Expected 0 commits in %s, got %d", r.Name, len(r.Commits))
		}
	}
}

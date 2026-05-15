package store

import (
	"sync"

	"gitstat/internal/model"
)

type Store struct {
	mu           sync.RWMutex
	repositories []model.Repository
}

var GlobalStore = &Store{}

func (s *Store) SetRepositories(repos []model.Repository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repositories = repos
}

func (s *Store) GetRepositories() []model.Repository {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repositories
}

func (s *Store) GetAllCommits() []model.Commit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allCommits []model.Commit
	for _, repo := range s.repositories {
		allCommits = append(allCommits, repo.Commits...)
	}
	return allCommits
}

package aggregator

import (
	"testing"
	"time"

	"gitstat/internal/model"
)

func TestAggregateOverview(t *testing.T) {
	commits := []model.Commit{
		{Author: "Alice", Additions: 10, Deletions: 5},
		{Author: "Bob", Additions: 20, Deletions: 10},
	}

	stats := AggregateOverview(commits, 1)

	if stats.TotalCommits != 2 {
		t.Errorf("Expected 2 commits, got %d", stats.TotalCommits)
	}

	if stats.ActiveAuthors != 2 {
		t.Errorf("Expected 2 authors, got %d", stats.ActiveAuthors)
	}

	if stats.TotalAdditions != 30 {
		t.Errorf("Expected 30 additions, got %d", stats.TotalAdditions)
	}
}

func TestAggregateByAuthor(t *testing.T) {
	commits := []model.Commit{
		{Author: "Alice", Additions: 10, Deletions: 5},
		{Author: "Alice", Additions: 20, Deletions: 10},
		{Author: "Bob", Additions: 5, Deletions: 2},
	}

	stats := AggregateByAuthor(commits)

	if len(stats) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(stats))
	}

	aliceStats := stats["Alice"]
	if aliceStats["commits"].(int) != 2 {
		t.Errorf("Expected Alice to have 2 commits, got %d", aliceStats["commits"])
	}
}

func TestAggregateByTime(t *testing.T) {
	now := time.Now()
	commits := []model.Commit{
		{Date: now, Additions: 10},
		{Date: now.AddDate(0, 0, 1), Additions: 20},
	}

	result := AggregateByTime(commits, "day")

	if len(result) != 2 {
		t.Errorf("Expected 2 days, got %d", len(result))
	}
}

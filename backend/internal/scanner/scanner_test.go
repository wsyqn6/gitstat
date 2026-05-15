package scanner

import (
	"testing"
	"time"
)

func TestCalculateCutoffTime(t *testing.T) {
	now := time.Now()

	cutoff7d := calculateCutoffTime("7d")
	expected7d := now.AddDate(0, 0, -7)

	if cutoff7d.Day() != expected7d.Day() {
		t.Errorf("7d cutoff mismatch: expected day %d, got %d", expected7d.Day(), cutoff7d.Day())
	}

	cutoff30d := calculateCutoffTime("30d")
	expected30d := now.AddDate(0, 0, -30)

	if cutoff30d.Day() != expected30d.Day() {
		t.Errorf("30d cutoff mismatch: expected day %d, got %d", expected30d.Day(), cutoff30d.Day())
	}

	cutoffAll := calculateCutoffTime("all")
	if !cutoffAll.IsZero() {
		t.Errorf("Expected zero time for 'all' range, got %v", cutoffAll)
	}

	cutoffDefault := calculateCutoffTime("")
	if !cutoffDefault.IsZero() {
		t.Errorf("Expected zero time for empty range, got %v", cutoffDefault)
	}
}

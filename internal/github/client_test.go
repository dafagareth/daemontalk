package github

import (
	"testing"
)

func TestFetchEmptyTokenAndUser(t *testing.T) {
	stats := Fetch("", "")
	if stats.Login != "" {
		t.Errorf("expected empty stats for empty user, got %q", stats.Login)
	}
}

func TestFetchContributionsNoToken(t *testing.T) {
	weeks, total := fetchContributions("octocat", "")
	if weeks != nil || total != 0 {
		t.Errorf("expected nil weeks and 0 total for empty token, got %v and %d", weeks, total)
	}
}

func TestStatsDataStructures(t *testing.T) {
	repo := Repo{
		Name:        "daemontalk",
		URL:         "https://github.com/dafagareth/daemontalk",
		Description: "DaemonTalk core",
		Stars:       42,
		Language:    "Go",
		Fork:        false,
	}

	if repo.Name != "daemontalk" || repo.Stars != 42 || repo.Fork {
		t.Errorf("unexpected Repo struct field values: %+v", repo)
	}

	day := ContribDay{Date: "2026-09-02", Count: 5}
	week := ContribWeek{Days: []ContribDay{day}}
	if len(week.Days) != 1 || week.Days[0].Count != 5 {
		t.Errorf("unexpected ContribWeek values: %+v", week)
	}
}

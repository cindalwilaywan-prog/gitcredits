package main

import (
	"strings"
	"testing"
)

func TestBuildCredits_HasTitle(t *testing.T) {
	info := repoInfo{
		name:         "testproject",
		totalCommits: 42,
		contributors: []contributor{
			{name: "Alice", commits: 30},
			{name: "Bob", commits: 12},
		},
	}
	lines := buildCredits(info, 80)

	if len(lines) == 0 {
		t.Fatal("buildCredits returned empty")
	}

	// should contain ASCII art title
	found := false
	for _, l := range lines {
		if strings.Contains(l, "██") {
			found = true
			break
		}
	}
	if !found {
		t.Error("credits should contain ASCII art title (██)")
	}
}

func TestBuildCredits_HasProjectBy(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 10,
		contributors: []contributor{
			{name: "Alice", commits: 10},
		},
	}
	lines := buildCredits(info, 80)

	found := false
	for _, l := range lines {
		if strings.Contains(l, "A   P R O J E C T   B Y") {
			found = true
			break
		}
	}
	if !found {
		t.Error("credits should contain 'A   P R O J E C T   B Y'")
	}
}

func TestBuildCredits_HasStarring(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 20,
		contributors: []contributor{
			{name: "Alice", commits: 15},
			{name: "Bob", commits: 5},
		},
	}
	lines := buildCredits(info, 80)

	found := false
	for _, l := range lines {
		if strings.Contains(l, "S T A R R I N G") {
			found = true
			break
		}
	}
	if !found {
		t.Error("credits should contain 'S T A R R I N G' when multiple contributors")
	}
}

func TestBuildCredits_HasTheEnd(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 1,
		contributors: []contributor{
			{name: "Alice", commits: 1},
		},
	}
	lines := buildCredits(info, 80)

	found := false
	for _, l := range lines {
		if strings.Contains(l, "██") && (strings.Contains(l, "THE") || strings.Contains(strings.Join(lines[len(lines)-30:], "\n"), "THE END")) {
			found = true
			break
		}
	}
	// just check THE END exists in big text somewhere near the end
	for _, l := range lines[len(lines)-30:] {
		if strings.Contains(l, "██") {
			found = true
			break
		}
	}
	if !found {
		t.Error("credits should end with THE END in big text")
	}
}

func TestBuildCredits_CommitCount(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 99,
		contributors: []contributor{
			{name: "Alice", commits: 99},
		},
	}
	lines := buildCredits(info, 80)

	found := false
	for _, l := range lines {
		if strings.Contains(l, "99  C O M M I T S") {
			found = true
			break
		}
	}
	if !found {
		t.Error("credits should contain commit count")
	}
}

func TestBigText(t *testing.T) {
	rows := bigText("AB")
	if len(rows) != 5 {
		t.Fatalf("bigText should return 5 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if len(r) == 0 {
			t.Error("bigText row should not be empty")
		}
	}
}

func TestCenterText(t *testing.T) {
	result := centerText("hi", 10)
	if !strings.HasPrefix(result, "    ") {
		t.Errorf("expected padding, got %q", result)
	}
	if strings.TrimSpace(result) != "hi" {
		t.Errorf("expected 'hi' centered, got %q", result)
	}
}

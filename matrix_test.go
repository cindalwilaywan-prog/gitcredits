package main

import (
	"strings"
	"testing"
)

func TestBuildMatrixCards_TitleCard(t *testing.T) {
	info := repoInfo{
		name:         "testproject",
		totalCommits: 50,
		contributors: []contributor{
			{name: "Alice", commits: 50},
		},
	}
	cards := buildMatrixCards(info, 80, 24)

	if len(cards) == 0 {
		t.Fatal("buildMatrixCards returned no cards")
	}

	// first card should have title
	found := false
	for _, l := range cards[0].lines {
		if strings.Contains(l, "██") {
			found = true
			break
		}
	}
	if !found {
		t.Error("first card should contain ASCII art title")
	}
}

func TestBuildMatrixCards_HeroCards(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 30,
		contributors: []contributor{
			{name: "Alice", commits: 20},
			{name: "Bob", commits: 10},
		},
	}
	cards := buildMatrixCards(info, 80, 24)

	// should have: title + 2 hero cards + stats + will return = 5 minimum
	if len(cards) < 4 {
		t.Errorf("expected at least 4 cards, got %d", len(cards))
	}
}

func TestBuildMatrixCards_WillReturn(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 1,
		contributors: []contributor{
			{name: "Alice", commits: 1},
		},
	}
	cards := buildMatrixCards(info, 80, 24)

	lastCard := cards[len(cards)-1]
	found := false
	for _, l := range lastCard.lines {
		if strings.Contains(l, "THE CONTRIBUTORS WILL RETURN") {
			found = true
			break
		}
	}
	if !found {
		t.Error("last card should contain 'THE CONTRIBUTORS WILL RETURN'")
	}
}

func TestMatrixHeroTitle(t *testing.T) {
	tests := []struct {
		rank    int
		commits int
		want    string
	}{
		{0, 200, "THE ARCHITECT"},
		{0, 50, "THE FOUNDER"},
		{1, 10, "THE GUARDIAN"},
	}

	for _, tt := range tests {
		got := matrixHeroTitle(tt.rank, tt.commits)
		if got != tt.want {
			t.Errorf("matrixHeroTitle(%d, %d) = %q, want %q", tt.rank, tt.commits, got, tt.want)
		}
	}
}

func TestBuildMatrixCards_CardHeight(t *testing.T) {
	info := repoInfo{
		name:         "test",
		totalCommits: 1,
		contributors: []contributor{
			{name: "Alice", commits: 1},
		},
	}
	cards := buildMatrixCards(info, 80, 24)

	for i, card := range cards {
		if len(card.lines) != 24 {
			t.Errorf("card %d has %d lines, want 24", i, len(card.lines))
		}
	}
}

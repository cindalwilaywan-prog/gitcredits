package main

import (
	"fmt"
	"strings"
)

// matrix rain characters
var matrixChars = []rune("ﾊﾐﾋｰｳｼﾅﾓﾆｻﾜﾂｵﾘｱﾎﾃﾏｹﾒｴｶｷﾑﾕﾗｾﾈｽﾀﾇﾍ012345789ABCDEFZ")

func matrixHeroTitle(rank int, commits int) string {
	if rank == 0 {
		if commits >= 100 {
			return "THE ARCHITECT"
		}
		return "THE FOUNDER"
	}
	titles := []string{
		"THE GUARDIAN",
		"THE WARRIOR",
		"THE SENTINEL",
		"THE VOYAGER",
		"THE PHOENIX",
		"THE TITAN",
		"THE RANGER",
		"THE STRIKER",
	}
	idx := rank - 1
	if idx >= len(titles) {
		idx = idx % len(titles)
	}
	return titles[idx]
}

// A "card" is one screen of content to display
type matrixCard struct {
	lines []string // text content, indexed by row (len = height)
}

func buildMatrixCards(info repoInfo, width, height int) []matrixCard {
	var cards []matrixCard

	center := func(s string) string {
		return centerText(s, width)
	}

	makeCard := func(content []string) matrixCard {
		lines := make([]string, height)
		startY := (height - len(content)) / 2
		if startY < 0 {
			startY = 0
		}
		for i, line := range content {
			if startY+i < height {
				lines[startY+i] = line
			}
		}
		return matrixCard{lines: lines}
	}

	// card 0: title (big + description + stats summary)
	var titleContent []string
	titleContent = append(titleContent, center("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	titleContent = append(titleContent, "")
	titleRows := bigText(info.name)
	for _, row := range titleRows {
		titleContent = append(titleContent, center(row))
	}
	titleContent = append(titleContent, "")
	if info.description != "" {
		titleContent = append(titleContent, center("\""+info.description+"\""))
		titleContent = append(titleContent, "")
	}
	if info.language != "" || info.stars > 0 {
		var meta []string
		if info.language != "" {
			meta = append(meta, info.language)
		}
		if info.stars > 0 {
			meta = append(meta, fmt.Sprintf("★ %d stars", info.stars))
		}
		if info.license != "" {
			meta = append(meta, info.license)
		}
		titleContent = append(titleContent, center(strings.Join(meta, "  ·  ")))
	}
	titleContent = append(titleContent, "")
	titleContent = append(titleContent, center("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	cards = append(cards, makeCard(titleContent))

	// hero cards
	spacedName := func(name string) string {
		upper := strings.ToUpper(name)
		runes := []rune(upper)
		var spaced []rune
		for i, r := range runes {
			spaced = append(spaced, r)
			if i < len(runes)-1 {
				if r == ' ' {
					spaced = append(spaced, ' ', ' ')
				} else {
					spaced = append(spaced, ' ')
				}
			}
		}
		return string(spaced)
	}

	for rank, c := range info.contributors {
		var content []string
		content = append(content, center("━━━━━━━━━━━━━━━━━━━━━━━━"))
		content = append(content, "")
		content = append(content, center(matrixHeroTitle(rank, c.commits)))
		content = append(content, "")
		content = append(content, "")
		content = append(content, center(spacedName(c.name)))
		content = append(content, "")
		content = append(content, center(fmt.Sprintf("⚡ %d commits ⚡", c.commits)))
		content = append(content, "")
		content = append(content, center("━━━━━━━━━━━━━━━━━━━━━━━━"))
		cards = append(cards, makeCard(content))
	}

	// highlights
	if len(info.highlights) > 0 {
		var content []string
		content = append(content, center("E P I C   M O M E N T S"))
		content = append(content, "")
		for _, h := range info.highlights {
			content = append(content, center("⚡ "+h))
		}
		cards = append(cards, makeCard(content))
	}

	// stats
	var statsContent []string
	statsContent = append(statsContent, center(fmt.Sprintf("%d COMMITS  ·  %d HEROES", info.totalCommits, len(info.contributors))))
	if info.stars > 0 {
		statsContent = append(statsContent, center(fmt.Sprintf("★ %d STARGAZERS ★", info.stars)))
	}
	if info.language != "" {
		statsContent = append(statsContent, center("Forged in "+info.language))
	}
	cards = append(cards, makeCard(statsContent))

	// will return
	cards = append(cards, makeCard([]string{
		center("THE CONTRIBUTORS WILL RETURN"),
	}))

	return cards
}

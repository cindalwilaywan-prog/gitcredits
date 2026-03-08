package main

import (
	"fmt"
	"strings"
)

func buildCredits(info repoInfo, width int) []string {
	var lines []string

	center := func(s string) string {
		return centerText(s, width)
	}

	blank := func(n int) {
		for i := 0; i < n; i++ {
			lines = append(lines, "")
		}
	}

	blank(20)

	titleRows := bigText(info.name)
	for _, row := range titleRows {
		lines = append(lines, center(row))
	}

	blank(2)

	if info.description != "" {
		lines = append(lines, center("\""+info.description+"\""))
	}

	blank(6)

	if len(info.contributors) > 0 {
		lines = append(lines, center("A   P R O J E C T   B Y"))
		blank(2)
		lines = append(lines, center(strings.ToUpper(info.contributors[0].name)))
		blank(1)
		lines = append(lines, center(fmt.Sprintf("— %d commits —", info.contributors[0].commits)))
	}

	blank(6)

	if len(info.contributors) > 1 {
		lines = append(lines, center("S T A R R I N G"))
		blank(2)
		for _, c := range info.contributors[1:] {
			lines = append(lines, center(strings.ToUpper(c.name)))
			lines = append(lines, center(fmt.Sprintf("%d commits", c.commits)))
			blank(1)
		}
	}

	blank(5)

	if len(info.highlights) > 0 {
		lines = append(lines, center("N O T A B L E   S C E N E S"))
		blank(2)
		for _, h := range info.highlights {
			lines = append(lines, center("· "+h+" ·"))
			blank(1)
		}
	}

	blank(5)

	lines = append(lines, center("━━━━━━━━━━━━━━━━━━━━"))
	blank(2)
	lines = append(lines, center(fmt.Sprintf("%d  C O M M I T S", info.totalCommits)))
	blank(1)
	lines = append(lines, center(fmt.Sprintf("%d  C O N T R I B U T O R S", len(info.contributors))))
	if info.stars > 0 {
		blank(1)
		lines = append(lines, center(fmt.Sprintf("★  %d  S T A R G A Z E R S  ★", info.stars)))
	}
	blank(2)
	if info.language != "" {
		lines = append(lines, center("Written in "+info.language))
	}
	if info.license != "" {
		lines = append(lines, center("Licensed under "+info.license))
	}
	blank(2)
	lines = append(lines, center("━━━━━━━━━━━━━━━━━━━━"))

	blank(6)

	endRows := bigText("THE END")
	for _, row := range endRows {
		lines = append(lines, center(row))
	}

	blank(20)

	return lines
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type repoInfo struct {
	name         string
	description  string
	totalCommits int
	contributors []contributor
	highlights   []string
	stars        int
	license      string
	language     string
}

type contributor struct {
	name    string
	commits int
}

func getRepoInfo() repoInfo {
	info := repoInfo{}

	if dir, err := os.Getwd(); err == nil {
		parts := strings.Split(dir, string(os.PathSeparator))
		info.name = parts[len(parts)-1]
	}

	if desc, err := os.ReadFile(".git/description"); err == nil {
		d := strings.TrimSpace(string(desc))
		if d != "" && d != "Unnamed repository; edit this file 'description' to name the repository." {
			info.description = d
		}
	}

	if info.description == "" {
		if out, err := exec.Command("gh", "repo", "view", "--json", "description", "-q", ".description").Output(); err == nil {
			d := strings.TrimSpace(string(out))
			if d != "" {
				info.description = d
			}
		}
	}

	if out, err := exec.Command("git", "rev-list", "--count", "HEAD").Output(); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.totalCommits = n
		}
	}

	if out, err := exec.Command("git", "shortlog", "-sn", "--no-merges", "HEAD").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				n, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
				info.contributors = append(info.contributors, contributor{
					name:    strings.TrimSpace(parts[1]),
					commits: n,
				})
			}
		}
	}

	sort.Slice(info.contributors, func(i, j int) bool {
		return info.contributors[i].commits > info.contributors[j].commits
	})

	if out, err := exec.Command("git", "log", "--oneline", "--no-merges", "-50", "--format=%s").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "feat:") || strings.HasPrefix(line, "fix:") {
				msg := line
				if strings.HasPrefix(line, "feat: ") {
					msg = line[6:]
				} else if strings.HasPrefix(line, "fix: ") {
					msg = line[5:]
				}
				info.highlights = append(info.highlights, msg)
				if len(info.highlights) >= 8 {
					break
				}
			}
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "stargazerCount", "-q", ".stargazerCount").Output(); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.stars = n
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "licenseInfo", "-q", ".licenseInfo.name").Output(); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.license = l
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "primaryLanguage", "-q", ".primaryLanguage.name").Output(); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.language = l
		}
	}

	return info
}

// simple block letter generator for title
func bigText(s string) []string {
	letters := map[rune][]string{
		'A': {"  ██  ", " █  █ ", " ████ ", " █  █ ", " █  █ "},
		'B': {" ███  ", " █  █ ", " ███  ", " █  █ ", " ███  "},
		'C': {"  ███ ", " █    ", " █    ", " █    ", "  ███ "},
		'D': {" ███  ", " █  █ ", " █  █ ", " █  █ ", " ███  "},
		'E': {" ████ ", " █    ", " ███  ", " █    ", " ████ "},
		'F': {" ████ ", " █    ", " ███  ", " █    ", " █    "},
		'G': {"  ███ ", " █    ", " █ ██ ", " █  █ ", "  ███ "},
		'H': {" █  █ ", " █  █ ", " ████ ", " █  █ ", " █  █ "},
		'I': {" ███ ", "  █  ", "  █  ", "  █  ", " ███ "},
		'J': {"  ███ ", "    █ ", "    █ ", " █  █ ", "  ██  "},
		'K': {" █  █ ", " █ █  ", " ██   ", " █ █  ", " █  █ "},
		'L': {" █    ", " █    ", " █    ", " █    ", " ████ "},
		'M': {" █   █ ", " ██ ██ ", " █ █ █ ", " █   █ ", " █   █ "},
		'N': {" █   █ ", " ██  █ ", " █ █ █ ", " █  ██ ", " █   █ "},
		'O': {"  ██  ", " █  █ ", " █  █ ", " █  █ ", "  ██  "},
		'P': {" ███  ", " █  █ ", " ███  ", " █    ", " █    "},
		'Q': {"  ██  ", " █  █ ", " █  █ ", " █ █  ", "  █ █ "},
		'R': {" ███  ", " █  █ ", " ███  ", " █ █  ", " █  █ "},
		'S': {"  ███ ", " █    ", "  ██  ", "    █ ", " ███  "},
		'T': {" █████ ", "   █   ", "   █   ", "   █   ", "   █   "},
		'U': {" █  █ ", " █  █ ", " █  █ ", " █  █ ", "  ██  "},
		'V': {" █  █ ", " █  █ ", " █  █ ", "  ██  ", "  ██  "},
		'W': {" █   █ ", " █   █ ", " █ █ █ ", " ██ ██ ", " █   █ "},
		'X': {" █  █ ", " █  █ ", "  ██  ", " █  █ ", " █  █ "},
		'Y': {" █  █ ", " █  █ ", "  ██  ", "  █   ", "  █   "},
		'Z': {" ████ ", "   █  ", "  █   ", " █    ", " ████ "},
		'-': {"      ", "      ", " ──── ", "      ", "      "},
		' ': {"   ", "   ", "   ", "   ", "   "},
		'_': {"      ", "      ", "      ", "      ", " ████ "},
	}

	upper := strings.ToUpper(s)
	rows := make([]string, 5)

	for _, ch := range upper {
		letter, ok := letters[ch]
		if !ok {
			letter = letters[' ']
		}
		for row := 0; row < 5; row++ {
			rows[row] += letter[row]
		}
	}

	return rows
}

func centerText(s string, width int) string {
	runeLen := len([]rune(s))
	if runeLen >= width {
		return s
	}
	pad := (width - runeLen) / 2
	return strings.Repeat(" ", pad) + s
}

func formatCommitCount(n int) string {
	return fmt.Sprintf("%d", n)
}

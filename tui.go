package main

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Matrix animation states
const (
	mvsRain     = 0 // pure rain
	mvsResolve  = 1 // rain chars → real text
	mvsShow     = 2 // text fully shown, dim rain bg
	mvsDissolve = 3 // text → rain chars
)

// Frame counts per state (at 50ms/frame = 20fps)
const (
	framesRain     = 30 // 1.5s
	framesResolve  = 25 // 1.25s
	framesShow     = 50 // 2.5s
	framesDissolve = 20 // 1s
)

// Rain column state
type rainColumn struct {
	headY   int
	speed   int // ticks per advance
	tickAcc int
	length  int // trail length
	active  bool
}

type tickMsg struct{}

type model struct {
	// default theme
	lines     []string
	offset    int
	starField starField

	// common
	height int
	width  int
	done   bool
	theme  string

	// matrix theme
	cards      []matrixCard
	cardIdx    int
	mState     int
	mFrame     int // frame counter within current state
	rainCols   []rainColumn
	rainGrid   [][]rune // [row][col] current rain characters
	resolveMap [][]bool // which cells have been resolved
}

func (m *model) initRain() {
	m.rainCols = make([]rainColumn, m.width)
	m.rainGrid = make([][]rune, m.height)
	m.resolveMap = make([][]bool, m.height)
	for r := 0; r < m.height; r++ {
		m.rainGrid[r] = make([]rune, m.width)
		m.resolveMap[r] = make([]bool, m.width)
	}
	// initialize some active columns
	for c := 0; c < m.width; c++ {
		if rand.Intn(3) == 0 {
			m.rainCols[c] = rainColumn{
				headY:  rand.Intn(m.height),
				speed:  1 + rand.Intn(3),
				length: 4 + rand.Intn(12),
				active: true,
			}
		}
	}
}

func (m *model) tickRain() {
	for c := 0; c < m.width; c++ {
		col := &m.rainCols[c]
		if !col.active {
			// randomly activate
			if rand.Intn(40) == 0 {
				col.headY = 0
				col.speed = 1 + rand.Intn(3)
				col.length = 4 + rand.Intn(12)
				col.active = true
				col.tickAcc = 0
			}
			continue
		}

		col.tickAcc++
		if col.tickAcc >= col.speed {
			col.tickAcc = 0
			col.headY++

			if col.headY-col.length > m.height {
				col.active = false
				continue
			}
		}

		// update rain grid for this column
		for r := 0; r < m.height; r++ {
			dist := col.headY - r
			if dist >= 0 && dist < col.length {
				m.rainGrid[r][c] = matrixChars[rand.Intn(len(matrixChars))]
			}
		}
	}
}

func (m *model) resetResolve() {
	for r := 0; r < m.height; r++ {
		for c := 0; c < m.width; c++ {
			m.resolveMap[r][c] = false
		}
	}
}

func (m model) Init() tea.Cmd {
	if m.theme == "matrix" {
		return tea.Tick(50*time.Millisecond, func(_ time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return tea.Tick(120*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit
		case "up":
			if m.theme != "matrix" {
				m.offset -= 3
				if m.offset < 0 {
					m.offset = 0
				}
			}
			return m, nil
		case "down":
			if m.theme != "matrix" {
				m.offset += 3
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tickMsg:
		if m.theme == "matrix" {
			return m.updateMatrix()
		}
		m.offset++
		if m.offset > len(m.lines) {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Tick(120*time.Millisecond, func(_ time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return m, nil
}

func (m model) updateMatrix() (tea.Model, tea.Cmd) {
	m.tickRain()
	m.mFrame++

	switch m.mState {
	case mvsRain:
		if m.mFrame >= framesRain {
			m.mState = mvsResolve
			m.mFrame = 0
			m.resetResolve()
		}
	case mvsResolve:
		// progressively resolve text cells
		if m.cardIdx < len(m.cards) {
			card := m.cards[m.cardIdx]
			progress := float64(m.mFrame) / float64(framesResolve)
			for r := 0; r < m.height; r++ {
				line := ""
				if r < len(card.lines) {
					line = card.lines[r]
				}
				runes := []rune(line)
				for c := 0; c < len(runes) && c < m.width; c++ {
					if runes[c] != ' ' && runes[c] != 0 && !m.resolveMap[r][c] {
						if rand.Float64() < progress*0.15 {
							m.resolveMap[r][c] = true
						}
					}
				}
			}
		}
		if m.mFrame >= framesResolve {
			// force all resolved
			if m.cardIdx < len(m.cards) {
				card := m.cards[m.cardIdx]
				for r := 0; r < m.height; r++ {
					line := ""
					if r < len(card.lines) {
						line = card.lines[r]
					}
					runes := []rune(line)
					for c := 0; c < len(runes) && c < m.width; c++ {
						if runes[c] != ' ' && runes[c] != 0 {
							m.resolveMap[r][c] = true
						}
					}
				}
			}
			m.mState = mvsShow
			m.mFrame = 0
		}
	case mvsShow:
		if m.mFrame >= framesShow {
			m.mState = mvsDissolve
			m.mFrame = 0
		}
	case mvsDissolve:
		// progressively un-resolve
		progress := float64(m.mFrame) / float64(framesDissolve)
		for r := 0; r < m.height; r++ {
			for c := 0; c < m.width; c++ {
				if m.resolveMap[r][c] && rand.Float64() < progress*0.15 {
					m.resolveMap[r][c] = false
				}
			}
		}
		if m.mFrame >= framesDissolve {
			m.cardIdx++
			if m.cardIdx >= len(m.cards) {
				m.done = true
				return m, tea.Quit
			}
			m.mState = mvsRain
			m.mFrame = 0
			m.resetResolve()
		}
	}

	return m, tea.Tick(50*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) View() string {
	if m.done {
		return ""
	}
	if m.theme == "matrix" {
		return m.viewMatrix()
	}
	return m.viewDefault()
}

func (m model) viewMatrix() string {
	// color palette
	greenBright := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF41"))
	greenMed := lipgloss.NewStyle().Foreground(lipgloss.Color("#00AA30"))
	greenDim := lipgloss.NewStyle().Foreground(lipgloss.Color("#005518"))
	greenVDim := lipgloss.NewStyle().Foreground(lipgloss.Color("#003310"))
	goldText := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	whiteText := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	cyanText := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))

	var card matrixCard
	if m.cardIdx < len(m.cards) {
		card = m.cards[m.cardIdx]
	}

	// find text block bounds
	textTop := m.height
	textBottom := 0
	if m.cardIdx < len(m.cards) {
		for r := 0; r < m.height; r++ {
			if r < len(card.lines) && strings.TrimSpace(card.lines[r]) != "" {
				if r < textTop {
					textTop = r
				}
				if r > textBottom {
					textBottom = r
				}
			}
		}
	}

	// clear box: centered rectangle with generous padding
	boxPadV := 3 // vertical padding
	boxTop := textTop - boxPadV
	boxBottom := textBottom + boxPadV
	if boxTop < 0 {
		boxTop = 0
	}
	if boxBottom >= m.height {
		boxBottom = m.height - 1
	}
	// horizontal: leave rain columns on edges (15% each side)
	rainEdge := m.width / 7
	boxLeft := rainEdge
	boxRight := m.width - rainEdge

	// is text currently visible?
	textVisible := m.mState == mvsResolve || m.mState == mvsShow || m.mState == mvsDissolve

	var sb strings.Builder

	for r := 0; r < m.height; r++ {
		for c := 0; c < m.width; c++ {
			// check if this cell has text
			isTextCell := false
			var textRune rune
			if m.cardIdx < len(m.cards) && r < len(card.lines) {
				runes := []rune(card.lines[r])
				if c < len(runes) && runes[c] != ' ' && runes[c] != 0 {
					isTextCell = true
					textRune = runes[c]
				}
			}

			resolved := r < len(m.resolveMap) && c < len(m.resolveMap[r]) && m.resolveMap[r][c]

			if isTextCell && resolved {
				// determine line content for color
				ch := string(textRune)
				lineStr := ""
				if r < len(card.lines) {
					lineStr = strings.TrimSpace(card.lines[r])
				}
				if strings.Contains(lineStr, "THE ") && !strings.Contains(lineStr, "WILL RETURN") && !strings.Contains(lineStr, "commits") {
					sb.WriteString(goldText.Render(ch))
				} else if strings.Contains(lineStr, "██") {
					sb.WriteString(goldText.Render(ch))
				} else if strings.Contains(lineStr, "\"") || strings.Contains(lineStr, "·") || strings.Contains(lineStr, "★") || strings.Contains(lineStr, "Forged") {
					sb.WriteString(cyanText.Render(ch))
				} else if strings.Contains(lineStr, "━") {
					sb.WriteString(goldText.Render(ch))
				} else {
					sb.WriteString(whiteText.Render(ch))
				}
			} else if isTextCell && m.mState == mvsResolve {
				// not yet resolved — show scrambled char
				ch := matrixChars[rand.Intn(len(matrixChars))]
				sb.WriteString(greenBright.Render(string(ch)))
			} else if textVisible && r >= boxTop && r <= boxBottom && c >= boxLeft && c <= boxRight {
				// inside clear box — black background
				sb.WriteRune(' ')
			} else {
				// rain background
				rainChar := m.rainGrid[r][c]
				if rainChar != 0 {
					col := m.rainCols[c]
					dist := col.headY - r
					if dist == 0 {
						sb.WriteString(greenBright.Render(string(rainChar)))
					} else if dist > 0 && dist < 3 {
						sb.WriteString(greenMed.Render(string(rainChar)))
					} else if dist >= 3 && dist < 6 {
						sb.WriteString(greenDim.Render(string(rainChar)))
					} else {
						sb.WriteString(greenVDim.Render(string(rainChar)))
					}
				} else {
					sb.WriteRune(' ')
				}
			}
		}
		if r < m.height-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

func (m model) viewDefault() string {
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Bold(true)
	silver := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	dimmer := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	scene := lipgloss.NewStyle().Foreground(lipgloss.Color("#B0C4DE"))
	starBright := lipgloss.NewStyle().Foreground(lipgloss.Color("#8899AA"))
	gold := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	contributor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	var screenLines []string
	start := m.offset
	end := m.offset + m.height

	if start < 0 {
		start = 0
	}

	for i := start; i < end && i < len(m.lines); i++ {
		line := m.lines[i]
		screenIdx := i - start
		trimmed := strings.TrimSpace(line)

		fadeTop := 4
		fadeBottom := 4
		distFromTop := screenIdx
		distFromBottom := m.height - 1 - screenIdx

		isFaded := distFromTop < fadeTop || distFromBottom < fadeBottom
		isVeryFaded := distFromTop < 2 || distFromBottom < 2

		var styled string
		if trimmed == "" {
			starLine := make([]rune, m.width)
			for j := range starLine {
				starLine[j] = ' '
			}
			for _, s := range m.starField.stars {
				if s.y == i && s.x < m.width {
					starLine[s.x] = s.ch
				}
			}
			sl := string(starLine)
			if strings.TrimSpace(sl) != "" {
				styled = starBright.Render(sl)
			} else {
				styled = ""
			}
		} else if strings.Contains(trimmed, "██") {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = title.Render(line)
			}
		} else if strings.Contains(trimmed, "━") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "★") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = gold.Render(line)
			}
		} else if strings.HasPrefix(trimmed, "A   P R O") || strings.HasPrefix(trimmed, "S T A R") ||
			strings.HasPrefix(trimmed, "N O T A B") {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = bright.Render(line)
			}
		} else if strings.Contains(trimmed, "C O M M") || strings.Contains(trimmed, "C O N T R") ||
			strings.Contains(trimmed, "S T A R G") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = silver.Render(line)
			}
		} else if strings.Contains(trimmed, "\"") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "· ") && strings.HasSuffix(trimmed, " ·") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = scene.Render(line)
			}
		} else if strings.Contains(trimmed, "—") && strings.Contains(trimmed, "commits") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "commits") && !strings.Contains(trimmed, "C O M") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if trimmed == strings.ToUpper(trimmed) && len(trimmed) > 2 && !strings.Contains(trimmed, " O ") {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = contributor.Render(line)
			}
		} else {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = silver.Render(line)
			}
		}

		screenLines = append(screenLines, styled)
	}

	for len(screenLines) < m.height {
		screenLines = append(screenLines, "")
	}

	return strings.Join(screenLines, "\n")
}

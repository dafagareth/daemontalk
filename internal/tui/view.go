package tui

import (
	"fmt"
	"strings"

	"daemontalk/internal/post"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.Ready {
		return "\n  Initializing DaemonTalk TUI..."
	}

	theme := m.GetTheme()

	// 2 lines reserved for bottom status bar
	panelHeight := m.Height - 2
	if panelHeight < 6 {
		panelHeight = 6
	}

	// Dynamic box styles based on active theme
	inactBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderInactive)

	actBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderActive)

	var panels string

	if m.IsFullReader {
		// Full Reader Mode
		scrollPercent := fmt.Sprintf("%.0f%%", m.Viewport.ScrollPercent()*100)
		if m.Viewport.AtTop() {
			scrollPercent = "TOP"
		} else if m.Viewport.AtBottom() {
			scrollPercent = "BOT"
		}

		var currentPost post.Post
		if it, ok := m.List.SelectedItem().(Item); ok {
			currentPost = it.Post
		} else if len(m.Posts) > 0 {
			currentPost = m.Posts[0]
		}

		fullTitle := fmt.Sprintf(" [Reading] %s (%s) ", currentPost.Title, scrollPercent)
		maxTitleLen := m.Width - 12
		runes := []rune(fullTitle)
		if len(runes) > maxTitleLen && maxTitleLen > 10 {
			fullTitle = string(runes[:maxTitleLen-4]) + "... (" + scrollPercent + ") "
		}
		fullTitleStyled := lipgloss.NewStyle().Foreground(theme.ActiveAccent).Bold(true).Render(fullTitle)

		panels = actBox.
			Width(m.Width - 2).
			Height(panelHeight - 2).
			Render(fullTitleStyled + "\n" + m.Viewport.View())
	} else {
		// Split Mode
		listWidth := m.Width * 35 / 100
		if listWidth < 34 {
			listWidth = 34
		}
		if listWidth > 52 {
			listWidth = 52
		}
		contentWidth := m.Width - listWidth

		lStyle := inactBox
		cStyle := inactBox
		if m.ActivePanel == 0 {
			lStyle = actBox
		} else {
			cStyle = actBox
		}

		// 1. Left Panel (Dispatches list)
		leftCount := fmt.Sprintf("%d/%d", m.List.Index()+1, len(m.Posts))
		if m.List.FilterValue() != "" {
			leftCount = fmt.Sprintf("filter: %q", m.List.FilterValue())
		}
		leftTitle := fmt.Sprintf(" [1] Dispatches (%s) ", leftCount)
		if m.ActivePanel == 0 {
			leftTitle = lipgloss.NewStyle().Foreground(theme.ActiveAccent).Bold(true).Render(leftTitle)
		} else {
			leftTitle = lipgloss.NewStyle().Foreground(theme.TextMuted).Render(leftTitle)
		}

		listStr := lStyle.
			Width(listWidth - 2).
			Height(panelHeight - 2).
			Render(leftTitle + "\n" + m.List.View())

		// 2. Right Panel (Article preview)
		scrollPercent := fmt.Sprintf("%.0f%%", m.Viewport.ScrollPercent()*100)
		if m.Viewport.AtTop() {
			scrollPercent = "TOP"
		} else if m.Viewport.AtBottom() {
			scrollPercent = "BOT"
		}
		rightTitle := fmt.Sprintf(" [2] Article Preview (%s) ", scrollPercent)
		if m.ActivePanel == 1 {
			rightTitle = lipgloss.NewStyle().Foreground(theme.ActiveAccent).Bold(true).Render(rightTitle)
		} else {
			rightTitle = lipgloss.NewStyle().Foreground(theme.TextMuted).Render(rightTitle)
		}

		contentStr := cStyle.
			Width(contentWidth - 2).
			Height(panelHeight - 2).
			Render(rightTitle + "\n" + m.Viewport.View())

		panels = lipgloss.JoinHorizontal(lipgloss.Top, listStr, contentStr)
	}

	// Dynamic status bar styling
	sBarStyle := lipgloss.NewStyle().
		Background(theme.StatusBg).
		Foreground(theme.TextNormal).
		Padding(0, 1)

	sKeyStyle := lipgloss.NewStyle().
		Background(theme.StatusBg).
		Foreground(theme.ActiveAccent).
		Bold(true)

	sDescStyle := lipgloss.NewStyle().
		Background(theme.StatusBg).
		Foreground(theme.TextMuted)

	// 3. Render 2-Line Bottom Status Bar
	var row1, row2 string

	if m.IsFullReader {
		var row1Items []string
		row1Items = append(row1Items, sKeyStyle.Render("Esc/Enter/q")+" "+sDescStyle.Render("Back to List"))
		row1Items = append(row1Items, sKeyStyle.Render("t")+" "+sDescStyle.Render("Theme: "+theme.Name))
		row1Items = append(row1Items, sKeyStyle.Render("↑/↓,j/k")+" "+sDescStyle.Render("Scroll"))
		row1Items = append(row1Items, sKeyStyle.Render("d/u")+" "+sDescStyle.Render("Half Page"))
		row1Items = append(row1Items, sKeyStyle.Render("g/G")+" "+sDescStyle.Render("Top/Bottom"))
		row1 = strings.Join(row1Items, "  │  ")

		if m.FlashMsg != "" {
			row2 = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render(m.FlashMsg)
		} else {
			var row2Items []string
			row2Items = append(row2Items, sKeyStyle.Render("o")+" "+sDescStyle.Render("Open Image"))
			row2Items = append(row2Items, sKeyStyle.Render("w")+" "+sDescStyle.Render("Open Article on Web"))
			row2 = strings.Join(row2Items, "  │  ")
		}
	} else {
		var row1Items []string
		row1Items = append(row1Items, sKeyStyle.Render("Enter")+" "+sDescStyle.Render("Full Reader"))
		row1Items = append(row1Items, sKeyStyle.Render("t")+" "+sDescStyle.Render("Theme: "+theme.Name))
		row1Items = append(row1Items, sKeyStyle.Render("Tab/1/2")+" "+sDescStyle.Render("Switch Panel"))
		row1Items = append(row1Items, sKeyStyle.Render("↑/↓,j/k")+" "+sDescStyle.Render("Navigate"))
		if m.ActivePanel == 0 {
			row1Items = append(row1Items, sKeyStyle.Render("/")+" "+sDescStyle.Render("Search"))
		} else {
			row1Items = append(row1Items, sKeyStyle.Render("d/u")+" "+sDescStyle.Render("Scroll Page"))
			row1Items = append(row1Items, sKeyStyle.Render("g/G")+" "+sDescStyle.Render("Top/Bottom"))
		}
		row1Items = append(row1Items, sKeyStyle.Render("q")+" "+sDescStyle.Render("Quit"))
		row1 = strings.Join(row1Items, "  │  ")

		if m.FlashMsg != "" {
			row2 = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render(m.FlashMsg)
		} else {
			var row2Items []string
			row2Items = append(row2Items, sKeyStyle.Render("o")+" "+sDescStyle.Render("Open Image"))
			row2Items = append(row2Items, sKeyStyle.Render("w")+" "+sDescStyle.Render("Open Article on Web"))
			row2 = strings.Join(row2Items, "  │  ")
		}
	}

	bar1 := sBarStyle.Width(m.Width).Render(" " + row1)
	bar2 := sBarStyle.Width(m.Width).Render(" " + row2)
	bar := bar1 + "\n" + bar2

	return lipgloss.JoinVertical(lipgloss.Left, panels, bar)
}

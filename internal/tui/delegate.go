package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"daemontalk/internal/post"
)

// Item wraps post.Post for bubbles/list.Item interface
type Item struct {
	Post post.Post
}

func (i Item) Title() string       { return i.Post.Title }
func (i Item) Description() string { return i.Post.Date.Format("02 Jan 2006") }
func (i Item) FilterValue() string { return i.Post.Title + " " + strings.Join(i.Post.Tags, " ") }

// LazyDelegate implements a compact, clean Lazygit-like row delegate
type LazyDelegate struct {
	Theme Theme
}

func (d LazyDelegate) Height() int                             { return 2 }
func (d LazyDelegate) Spacing() int                            { return 0 }
func (d LazyDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d LazyDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(Item)
	if !ok {
		return
	}

	theme := d.Theme
	if theme.Name == "" {
		theme = Themes[0]
	}
	p := it.Post
	isSelected := index == m.Index()
	width := m.Width()

	maxTitleWidth := width - 4
	if maxTitleWidth < 5 {
		maxTitleWidth = 5
	}

	// Title truncation (visual column width safe)
	title := p.Title
	if lipgloss.Width(title) > maxTitleWidth {
		for len(title) > 0 && lipgloss.Width(title+"…") > maxTitleWidth {
			runes := []rune(title)
			title = string(runes[:len(runes)-1])
		}
		title += "…"
	}

	// Subtitle (Date + Tags, visual column width safe)
	dateStr := p.Date.Format("02 Jan 2006")
	tagStr := ""
	if len(p.Tags) > 0 {
		tagStr = " • " + strings.Join(p.Tags, ", ")
	}
	sub := dateStr + tagStr
	if lipgloss.Width(sub) > maxTitleWidth {
		for len(sub) > 0 && lipgloss.Width(sub+"…") > maxTitleWidth {
			runes := []rune(sub)
			sub = string(runes[:len(runes)-1])
		}
		sub += "…"
	}

	// Pad with spaces to fill exact visual column width
	titlePad := ""
	if rem := maxTitleWidth - lipgloss.Width(title); rem > 0 {
		titlePad = strings.Repeat(" ", rem)
	}
	subPad := ""
	if rem := maxTitleWidth - lipgloss.Width(sub); rem > 0 {
		subPad = strings.Repeat(" ", rem)
	}

	var titleLine, subLine string
	if isSelected {
		cursor := lipgloss.NewStyle().Foreground(theme.ActiveAccent).Bold(true).Render("▎ ")
		titleStyled := lipgloss.NewStyle().
			Foreground(theme.TextNormal).
			Background(theme.SelectedBg).
			Bold(true).
			Render(title + titlePad)
		subStyled := lipgloss.NewStyle().
			Foreground(theme.ActiveAccent).
			Background(theme.SelectedBg).
			Render(sub + subPad)

		titleLine = cursor + titleStyled
		subLine = "  " + subStyled
	} else {
		titleStyled := lipgloss.NewStyle().
			Foreground(theme.TextNormal).
			Render(title)
		subStyled := lipgloss.NewStyle().
			Foreground(theme.TextMuted).
			Render(sub)

		titleLine = "  " + titleStyled
		subLine = "  " + subStyled
	}

	fmt.Fprintf(w, "%s\n%s", titleLine, subLine)
}

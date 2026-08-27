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

	// Title truncation (rune-safe for multi-byte Unicode/emojis)
	title := p.Title
	maxTitleWidth := width - 4
	if maxTitleWidth < 5 {
		maxTitleWidth = 5
	}
	titleRunes := []rune(title)
	if len(titleRunes) > maxTitleWidth {
		title = string(titleRunes[:maxTitleWidth-1]) + "…"
	}

	// Subtitle (Date + Tags, rune-safe)
	dateStr := p.Date.Format("02 Jan 2006")
	tagStr := ""
	if len(p.Tags) > 0 {
		tagStr = " • " + strings.Join(p.Tags, ", ")
	}
	sub := dateStr + tagStr
	subRunes := []rune(sub)
	if len(subRunes) > maxTitleWidth {
		sub = string(subRunes[:maxTitleWidth-1]) + "…"
	}

	var titleLine, subLine string
	if isSelected {
		cursor := lipgloss.NewStyle().Foreground(theme.ActiveAccent).Bold(true).Render("▎ ")
		titleStyled := lipgloss.NewStyle().
			Foreground(theme.TextNormal).
			Background(theme.SelectedBg).
			Bold(true).
			Render(fmt.Sprintf("%-*s", maxTitleWidth, title))
		subStyled := lipgloss.NewStyle().
			Foreground(theme.ActiveAccent).
			Background(theme.SelectedBg).
			Render(fmt.Sprintf("%-*s", maxTitleWidth, sub))

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

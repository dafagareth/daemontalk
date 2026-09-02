package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"daemontalk/internal/post"
)

type Model struct {
	List         list.Model
	Viewport     viewport.Model
	Posts        []post.Post
	ActivePanel  int // 0 = list, 1 = viewport
	IsFullReader bool
	IsCompact    bool // true when terminal width < 72
	ThemeIdx     int
	FlashMsg     string
	Ready        bool
	Width        int
	Height       int
}

func (m Model) GetTheme() Theme {
	if m.ThemeIdx < 0 || m.ThemeIdx >= len(Themes) {
		return Themes[0]
	}
	return Themes[m.ThemeIdx]
}

func NewModel() Model {
	postsDir := getContentPostsDir()
	posts, err := post.LoadAll(postsDir)
	if err != nil {
		fmt.Printf("Error loading posts from %s: %v\n", postsDir, err)
		posts = nil
	}

	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = Item{Post: p}
	}

	l := list.New(items, LazyDelegate{Theme: Themes[0]}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()

	return Model{
		List:         l,
		Posts:        posts,
		ActivePanel:  0,
		IsFullReader: false,
		ThemeIdx:     0, // 0 = Nord default
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) RenderCurrentPost() string {
	if len(m.Posts) == 0 {
		return "\n  No dispatches found in content/posts directory."
	}

	if it, ok := m.List.SelectedItem().(Item); ok {
		wrapWidth := m.Viewport.Width - 2
		if wrapWidth < 20 {
			wrapWidth = 20
		}
		return RenderPostMarkdown(it.Post, wrapWidth, m.GetTheme())
	}

	if m.List.FilterValue() != "" {
		return fmt.Sprintf("\n  No dispatches match filter: %q\n\n  Press [Esc] to clear search.", m.List.FilterValue())
	}

	idx := m.List.Index()
	if idx >= 0 && idx < len(m.Posts) {
		wrapWidth := m.Viewport.Width - 2
		if wrapWidth < 20 {
			wrapWidth = 20
		}
		return RenderPostMarkdown(m.Posts[idx], wrapWidth, m.GetTheme())
	}

	return "\n  Select a dispatch to read."
}

func (m *Model) RecalcSizes() {
	if m.Width <= 0 || m.Height <= 0 {
		return
	}

	m.IsCompact = m.Width < 72

	// Reserve 2 lines for the 2-row bottom status bar
	panelHeight := m.Height - 2
	if panelHeight < 6 {
		panelHeight = 6
	}

	if m.IsFullReader || (m.IsCompact && m.ActivePanel == 1) {
		contentInnerW := m.Width - 2
		if contentInnerW < 10 {
			contentInnerW = 10
		}
		contentInnerH := panelHeight - 2
		if contentInnerH < 4 {
			contentInnerH = 4
		}
		m.Viewport.Width = contentInnerW - 2
		m.Viewport.Height = contentInnerH - 1
		m.Viewport.SetContent(m.RenderCurrentPost())
	} else if m.IsCompact && m.ActivePanel == 0 {
		innerWidth := m.Width - 2
		if innerWidth < 10 {
			innerWidth = 10
		}
		innerHeight := panelHeight - 2
		if innerHeight < 4 {
			innerHeight = 4
		}
		m.List.SetSize(innerWidth, innerHeight-1)
	} else {
		listWidth := m.Width * 35 / 100
		if listWidth < 30 {
			listWidth = 30
		}
		if listWidth > 50 {
			listWidth = 50
		}
		contentWidth := m.Width - listWidth
		if contentWidth < 30 {
			contentWidth = 30
		}

		innerWidth := listWidth - 2
		innerHeight := panelHeight - 2

		m.List.SetSize(innerWidth, innerHeight-1)

		contentInnerW := contentWidth - 2
		contentInnerH := panelHeight - 2
		m.Viewport.Width = contentInnerW - 2
		m.Viewport.Height = contentInnerH - 1
		m.Viewport.SetContent(m.RenderCurrentPost())
	}
}

package tui

import (
	"fmt"
	"os"

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
	posts, err := post.LoadAll("content/posts")
	if err != nil {
		fmt.Printf("Error loading posts: %v\n", err)
		os.Exit(1)
	}

	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = Item{Post: p}
	}

	l := list.New(items, LazyDelegate{}, 0, 0)
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
		return "No posts found."
	}

	idx := m.List.Index()
	if idx < 0 || idx >= len(m.Posts) {
		return ""
	}

	wrapWidth := m.Viewport.Width - 4
	return RenderPostMarkdown(m.Posts[idx], wrapWidth, m.GetTheme())
}

func (m *Model) RecalcSizes() {
	if m.Width <= 0 || m.Height <= 0 {
		return
	}

	// Reserve 2 lines for the 2-row bottom status bar
	panelHeight := m.Height - 2
	if panelHeight < 6 {
		panelHeight = 6
	}

	if m.IsFullReader {
		contentInnerW := m.Width - 2
		contentInnerH := panelHeight - 2
		m.Viewport.Width = contentInnerW - 4
		m.Viewport.Height = contentInnerH - 1
		m.Viewport.SetContent(m.RenderCurrentPost())
	} else {
		listWidth := m.Width * 35 / 100
		if listWidth < 34 {
			listWidth = 34
		}
		if listWidth > 52 {
			listWidth = 52
		}
		contentWidth := m.Width - listWidth

		innerWidth := listWidth - 2
		innerHeight := panelHeight - 2

		// List size: subtract 1 from height for the custom title line
		m.List.SetSize(innerWidth, innerHeight-1)

		contentInnerW := contentWidth - 2
		contentInnerH := panelHeight - 2
		m.Viewport.Width = contentInnerW - 4
		m.Viewport.Height = contentInnerH - 1
		m.Viewport.SetContent(m.RenderCurrentPost())
	}
}

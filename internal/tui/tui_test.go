package tui

import (
	"bytes"
	"testing"
	"time"

	"daemontalk/internal/post"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m.ThemeIdx != 0 {
		t.Errorf("expected default ThemeIdx 0, got %d", m.ThemeIdx)
	}
	theme := m.GetTheme()
	if theme.Name != "Nord" {
		t.Errorf("expected Nord theme, got %s", theme.Name)
	}
}

func TestRecalcSizes_VariousDimensions(t *testing.T) {
	m := NewModel()

	tests := []struct {
		name    string
		width   int
		height  int
		compact bool
	}{
		{"Zero size", 0, 0, false},
		{"Mobile SSH / Small terminal", 40, 20, true},
		{"Compact boundary", 71, 24, true},
		{"Standard 80x24", 80, 24, false},
		{"Large 160x50", 160, 50, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m.Width = tc.width
			m.Height = tc.height
			m.RecalcSizes()

			if tc.width > 0 && tc.compact != m.IsCompact {
				t.Errorf("expected IsCompact=%v for width=%d, got %v", tc.compact, tc.width, m.IsCompact)
			}
		})
	}
}

func TestThemeCycling(t *testing.T) {
	m := NewModel()
	initialTheme := m.GetTheme().Name

	// Cycle theme with 't'
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model := newM.(Model)

	if model.GetTheme().Name == initialTheme {
		t.Errorf("expected theme to cycle from %s", initialTheme)
	}

	// Cycle 6 more times to complete loop of 7 themes
	for i := 0; i < len(Themes)-1; i++ {
		newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
		model = newM.(Model)
	}

	if model.GetTheme().Name != initialTheme {
		t.Errorf("expected theme to return to %s, got %s", initialTheme, model.GetTheme().Name)
	}
}

func TestDelegateRender(t *testing.T) {
	delegate := LazyDelegate{Theme: Themes[0]}
	p := post.Post{
		Title:       "Uji Coba TUI dengan Unicode — 🚀",
		Description: "Deskripsi singkat pengujian.",
		Date:        time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC),
		Tags:        []string{"tui", "charm", "golang"},
	}
	item := Item{Post: p}

	items := []list.Item{item}
	l := list.New(items, delegate, 60, 20)

	var buf bytes.Buffer
	delegate.Render(&buf, l, 0, item)

	out := buf.String()
	if out == "" {
		t.Errorf("expected non-empty rendered string from LazyDelegate")
	}
}

func TestNavigationUpdates(t *testing.T) {
	m := NewModel()
	m.Width = 100
	m.Height = 30
	m.Ready = true
	m.RecalcSizes()

	// Tab should switch panel to 1
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := newM.(Model)
	if model.ActivePanel != 1 {
		t.Errorf("expected ActivePanel 1, got %d", model.ActivePanel)
	}

	// Enter should toggle full reader mode
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newM.(Model)
	if !model.IsFullReader {
		t.Errorf("expected IsFullReader true")
	}

	// Esc should exit full reader mode
	newM, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = newM.(Model)
	if model.IsFullReader {
		t.Errorf("expected IsFullReader false after Esc")
	}
}

func TestRenderPostMarkdown(t *testing.T) {
	p := post.Post{
		Title:       "Test Markdown Rendering",
		Description: "Testing markdown rendering for TUI.",
		Date:        time.Now(),
		ReadTime:    3,
		Tags:        []string{"test"},
	}

	out := RenderPostMarkdown(p, 60, Themes[0])
	if out == "" {
		t.Fatalf("expected non-empty output from RenderPostMarkdown")
	}

	// Ensure cached result is returned quickly on second call
	outCached := RenderPostMarkdown(p, 60, Themes[0])
	if outCached != out {
		t.Errorf("expected cached output to match original output")
	}
}

func TestResolvePostURL(t *testing.T) {
	// External URL
	external := "https://images.unsplash.com/photo-1591799264318-7e6ef8ddb7ea?auto=format&fit=crop&w=1200&q=80"
	if got := ResolvePostURL(external, "slug"); got != external {
		t.Errorf("expected %s, got %s", external, got)
	}

	// Relative static image path
	rel := "/static/images/posts/arch.png"
	gotRel := ResolvePostURL(rel, "slug")
	if gotRel != "https://www.daemontalk.com/static/images/posts/arch.png" {
		t.Errorf("expected full static URL, got %s", gotRel)
	}

	// Fallback to slug
	gotSlug := ResolvePostURL("", "my-post")
	if gotSlug != "https://www.daemontalk.com/blog/my-post" {
		t.Errorf("expected blog slug URL, got %s", gotSlug)
	}
}

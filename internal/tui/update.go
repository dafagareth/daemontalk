package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Reset flash message on keypress
		m.FlashMsg = ""

		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// When searching in list, let list handle all typing
		if !m.IsFullReader && m.List.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "t", "T":
			// Cycle through the 7 themes
			m.ThemeIdx = (m.ThemeIdx + 1) % len(Themes)
			ActiveThemeIdx = m.ThemeIdx
			theme := m.GetTheme()
			m.FlashMsg = fmt.Sprintf("✓ Theme: %s (%d/%d)", theme.Name, m.ThemeIdx+1, len(Themes))
			m.Viewport.SetContent(m.RenderCurrentPost())
			return m, nil

		case "o":
			// Open cover image in browser (or article if no cover)
			if len(m.Posts) > 0 {
				p := m.Posts[m.List.Index()]
				baseURL := os.Getenv("BASE_URL")
				if baseURL == "" {
					baseURL = "https://daemontalk.com"
				}
				baseURL = strings.TrimSuffix(baseURL, "/")

				var targetURL string
				if p.Cover != "" {
					coverPath := p.Cover
					if !strings.HasPrefix(coverPath, "/") {
						coverPath = "/" + coverPath
					}
					targetURL = fmt.Sprintf("%s%s", baseURL, coverPath)
				} else {
					targetURL = fmt.Sprintf("%s/blog/%s", baseURL, p.Slug)
				}
				_ = openInBrowser(targetURL)
				clickable := OSC8Link(targetURL, targetURL)
				m.FlashMsg = fmt.Sprintf("✓ Copied to clipboard & Clickable: %s", clickable)
				return m, tea.Printf("%s", OSC52Copy(targetURL))
			}
			return m, nil

		case "w":
			// Open full article in browser
			if len(m.Posts) > 0 {
				p := m.Posts[m.List.Index()]
				baseURL := os.Getenv("BASE_URL")
				if baseURL == "" {
					baseURL = "https://daemontalk.com"
				}
				baseURL = strings.TrimSuffix(baseURL, "/")

				targetURL := fmt.Sprintf("%s/blog/%s", baseURL, p.Slug)
				_ = openInBrowser(targetURL)
				clickable := OSC8Link(targetURL, targetURL)
				m.FlashMsg = fmt.Sprintf("✓ Copied to clipboard & Clickable: %s", clickable)
				return m, tea.Printf("%s", OSC52Copy(targetURL))
			}
			return m, nil

		case "enter":
			if !m.IsFullReader {
				m.IsFullReader = true
				m.RecalcSizes()
				m.Viewport.GotoTop()
				return m, nil
			}
			m.IsFullReader = false
			m.RecalcSizes()
			return m, nil

		case "esc", "backspace":
			if m.IsFullReader {
				m.IsFullReader = false
				m.RecalcSizes()
				return m, nil
			}

		case "q":
			if m.IsFullReader {
				m.IsFullReader = false
				m.RecalcSizes()
				return m, nil
			}
			if m.ActivePanel == 1 {
				m.ActivePanel = 0
				return m, nil
			}
			return m, tea.Quit

		case "tab":
			if !m.IsFullReader {
				m.ActivePanel = (m.ActivePanel + 1) % 2
				return m, nil
			}

		case "1":
			if !m.IsFullReader {
				m.ActivePanel = 0
				return m, nil
			}

		case "2":
			if !m.IsFullReader {
				m.ActivePanel = 1
				return m, nil
			}

		case "g":
			if m.IsFullReader || m.ActivePanel == 1 {
				m.Viewport.GotoTop()
				return m, nil
			}

		case "G":
			if m.IsFullReader || m.ActivePanel == 1 {
				m.Viewport.GotoBottom()
				return m, nil
			}
		}

		// Route input to the active panel
		if !m.IsFullReader && m.ActivePanel == 0 {
			prevIdx := m.List.Index()
			var cmd tea.Cmd
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)

			// If selection changed, update right panel
			if m.List.Index() != prevIdx && m.Ready {
				m.Viewport.SetContent(m.RenderCurrentPost())
				m.Viewport.GotoTop()
			}
		} else {
			var cmd tea.Cmd
			m.Viewport, cmd = m.Viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true
		ActiveThemeIdx = m.ThemeIdx
		m.RecalcSizes()
	}

	return m, tea.Batch(cmds...)
}

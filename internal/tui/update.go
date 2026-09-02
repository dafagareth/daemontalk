package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Mouse wheel and click handling
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			if m.IsFullReader || m.ActivePanel == 1 || (m.IsCompact && m.ActivePanel == 1) {
				var cmd tea.Cmd
				m.Viewport, cmd = m.Viewport.Update(msg)
				return m, cmd
			} else {
				prevIdx := m.List.Index()
				var cmd tea.Cmd
				m.List, cmd = m.List.Update(msg)
				if m.List.Index() != prevIdx {
					m.Viewport.SetContent(m.RenderCurrentPost())
					m.Viewport.GotoTop()
				}
				return m, cmd
			}
		}

	case tea.KeyMsg:
		// Reset flash message on next keypress
		m.FlashMsg = ""

		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// When searching in list, synchronize preview in real time
		if !m.IsFullReader && m.List.FilterState() == list.Filtering {
			prevFilter := m.List.FilterValue()
			prevIdx := m.List.Index()
			var cmd tea.Cmd
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)

			if m.List.FilterValue() != prevFilter || m.List.Index() != prevIdx {
				m.Viewport.SetContent(m.RenderCurrentPost())
				m.Viewport.GotoTop()
			}
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "t", "T":
			// Cycle through the 7 developer themes
			m.ThemeIdx = (m.ThemeIdx + 1) % len(Themes)
			theme := m.GetTheme()
			m.List.SetDelegate(LazyDelegate{Theme: theme})
			m.FlashMsg = fmt.Sprintf("Theme: %s (%d/%d)", theme.Name, m.ThemeIdx+1, len(Themes))
			m.Viewport.SetContent(m.RenderCurrentPost())
			return m, nil

		case "o":
			// Open cover image / link
			if it, ok := m.List.SelectedItem().(Item); ok {
				p := it.Post
				targetURL := ResolvePostURL(p.Cover, p.Slug)
				_ = openInBrowser(targetURL)
				clickable := OSC8Link(targetURL, targetURL)
				m.FlashMsg = fmt.Sprintf("Copied link (OSC52): %s", clickable)
				return m, nil
			}
			return m, nil

		case "w":
			// Open article on web
			if it, ok := m.List.SelectedItem().(Item); ok {
				p := it.Post
				targetURL := ResolvePostURL("", p.Slug)
				_ = openInBrowser(targetURL)
				clickable := OSC8Link(targetURL, targetURL)
				m.FlashMsg = fmt.Sprintf("Copied web link: %s", clickable)
				return m, nil
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

		case "esc":
			if m.IsFullReader {
				m.IsFullReader = false
				m.RecalcSizes()
				return m, nil
			}
			if m.List.FilterState() == list.FilterApplied {
				m.List.ResetFilter()
				m.Viewport.SetContent(m.RenderCurrentPost())
				return m, nil
			}
			if m.ActivePanel == 1 {
				m.ActivePanel = 0
				return m, nil
			}

		case "backspace":
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
				if m.IsCompact {
					m.RecalcSizes()
				}
				return m, nil
			}

		case "h", "left":
			if !m.IsFullReader && m.ActivePanel == 1 {
				m.ActivePanel = 0
				if m.IsCompact {
					m.RecalcSizes()
				}
				return m, nil
			}

		case "l", "right":
			if !m.IsFullReader && m.ActivePanel == 0 {
				m.ActivePanel = 1
				if m.IsCompact {
					m.RecalcSizes()
				}
				return m, nil
			}

		case "1":
			if !m.IsFullReader {
				m.ActivePanel = 0
				if m.IsCompact {
					m.RecalcSizes()
				}
				return m, nil
			}

		case "2":
			if !m.IsFullReader {
				m.ActivePanel = 1
				if m.IsCompact {
					m.RecalcSizes()
				}
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

		case "d", "ctrl+d":
			if m.IsFullReader || m.ActivePanel == 1 {
				m.Viewport.HalfViewDown()
				return m, nil
			}

		case "u", "ctrl+u":
			if m.IsFullReader || m.ActivePanel == 1 {
				m.Viewport.HalfViewUp()
				return m, nil
			}

		case " ", "space":
			if m.IsFullReader || m.ActivePanel == 1 {
				m.Viewport.ViewDown()
				return m, nil
			}

		case "?":
			m.FlashMsg = "Shortcuts: [Enter] Full Reader · [Tab/1/2] Panels · [/] Search · [t] Theme · [o] Image · [w] Web · [q] Quit"
			return m, nil
		}

		// Route navigation input to the active panel
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
		m.RecalcSizes()
	}

	return m, tea.Batch(cmds...)
}

package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines the complete color palette and style for the TUI
type Theme struct {
	Name           string
	BorderInactive lipgloss.Color
	BorderActive   lipgloss.Color
	ActiveAccent   lipgloss.Color
	SelectedBg     lipgloss.Color
	TextNormal     lipgloss.Color
	TextMuted      lipgloss.Color
	StatusBg       lipgloss.Color
	GlamourStyle   string
}

// 7 Developer Color Themes
var Themes = []Theme{
	{
		Name:           "Nord",
		BorderInactive: lipgloss.Color("#3b4252"),
		BorderActive:   lipgloss.Color("#88c0d0"),
		ActiveAccent:   lipgloss.Color("#88c0d0"),
		SelectedBg:     lipgloss.Color("#2e3440"),
		TextNormal:     lipgloss.Color("#eceff4"),
		TextMuted:      lipgloss.Color("#616e88"),
		StatusBg:       lipgloss.Color("#2e3440"),
		GlamourStyle:   "dark",
	},
	{
		Name:           "Tokyo Night",
		BorderInactive: lipgloss.Color("#292e42"),
		BorderActive:   lipgloss.Color("#7aa2f7"),
		ActiveAccent:   lipgloss.Color("#bb9af7"),
		SelectedBg:     lipgloss.Color("#1f2335"),
		TextNormal:     lipgloss.Color("#c0caf5"),
		TextMuted:      lipgloss.Color("#565f89"),
		StatusBg:       lipgloss.Color("#1a1b26"),
		GlamourStyle:   "tokyo-night",
	},
	{
		Name:           "Catppuccin Mocha",
		BorderInactive: lipgloss.Color("#313244"),
		BorderActive:   lipgloss.Color("#cba6f7"),
		ActiveAccent:   lipgloss.Color("#89b4fa"),
		SelectedBg:     lipgloss.Color("#1e1e2e"),
		TextNormal:     lipgloss.Color("#cdd6f4"),
		TextMuted:      lipgloss.Color("#6c7086"),
		StatusBg:       lipgloss.Color("#181825"),
		GlamourStyle:   "dark",
	},
	{
		Name:           "Gruvbox Dark",
		BorderInactive: lipgloss.Color("#3c3836"),
		BorderActive:   lipgloss.Color("#fabd2f"),
		ActiveAccent:   lipgloss.Color("#fe8019"),
		SelectedBg:     lipgloss.Color("#32302f"),
		TextNormal:     lipgloss.Color("#ebdbb2"),
		TextMuted:      lipgloss.Color("#928374"),
		StatusBg:       lipgloss.Color("#282828"),
		GlamourStyle:   "dark",
	},
	{
		Name:           "Dracula",
		BorderInactive: lipgloss.Color("#44475a"),
		BorderActive:   lipgloss.Color("#bd93f9"),
		ActiveAccent:   lipgloss.Color("#ff79c6"),
		SelectedBg:     lipgloss.Color("#343746"),
		TextNormal:     lipgloss.Color("#f8f8f2"),
		TextMuted:      lipgloss.Color("#6272a4"),
		StatusBg:       lipgloss.Color("#282a36"),
		GlamourStyle:   "dracula",
	},
	{
		Name:           "Rose Pine",
		BorderInactive: lipgloss.Color("#26233a"),
		BorderActive:   lipgloss.Color("#ebbcba"),
		ActiveAccent:   lipgloss.Color("#31748f"),
		SelectedBg:     lipgloss.Color("#1f1d2e"),
		TextNormal:     lipgloss.Color("#e0def4"),
		TextMuted:      lipgloss.Color("#6e6a86"),
		StatusBg:       lipgloss.Color("#191724"),
		GlamourStyle:   "dark",
	},
	{
		Name:           "Monokai Pro",
		BorderInactive: lipgloss.Color("#3a3839"),
		BorderActive:   lipgloss.Color("#ffd866"),
		ActiveAccent:   lipgloss.Color("#a9dc76"),
		SelectedBg:     lipgloss.Color("#2d2a2e"),
		TextNormal:     lipgloss.Color("#fcfcfa"),
		TextMuted:      lipgloss.Color("#727072"),
		StatusBg:       lipgloss.Color("#221f22"),
		GlamourStyle:   "dark",
	},
}

var ActiveThemeIdx = 0

func GetActiveTheme() Theme {
	if ActiveThemeIdx < 0 || ActiveThemeIdx >= len(Themes) {
		return Themes[0]
	}
	return Themes[ActiveThemeIdx]
}

var ColorSuccess = lipgloss.Color("#a3be8c")

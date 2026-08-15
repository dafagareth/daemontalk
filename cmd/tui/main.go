package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"daemontalk/internal/tui"
)

func main() {
	program := tea.NewProgram(
		tui.NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting DaemonTalk TUI: %v\n", err)
		os.Exit(1)
	}
}

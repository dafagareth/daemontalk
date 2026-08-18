package tuisrv

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"daemontalk/internal/tui"
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal PTY")
		return nil, nil
	}

	renderer := bubbletea.MakeRenderer(s)
	lipgloss.SetColorProfile(renderer.ColorProfile())

	m := tui.NewModel()
	m.Width = pty.Window.Width
	m.Height = pty.Window.Height
	m.Ready = true
	m.RecalcSizes()

	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

// Start launches the Wish SSH server on the given address
func Start(addr string, hostKeyPath string) (*ssh.Server, error) {
	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithIdleTimeout(30*time.Minute),
		wish.WithMaxTimeout(2*time.Hour),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create wish ssh server: %w", err)
	}

	return s, nil
}

package tui

import (
	"os/exec"
	"regexp"
	"runtime"
)

var (
	reMarkdownImage = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

// openInBrowser opens a specified URL or file in the default system browser
func openInBrowser(target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	return cmd.Start()
}

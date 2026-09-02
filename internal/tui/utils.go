package tui

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var (
	reMarkdownImage = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	reBlockMath     = regexp.MustCompile(`(?s)\$\$(.*?)\$\$`)
	reInlineMath    = regexp.MustCompile(`\$([^$\n]+)\$`)
	reHTMLTags      = regexp.MustCompile(`<[^>]+>`)
	reCalloutOpen   = regexp.MustCompile(`(?i)<callout\s+type="([^"]+)">`)
)

// ResolvePostURL resolves a path or slug into a full URL (handling already-absolute http/https URLs)
func ResolvePostURL(rawPath, fallbackSlug string) string {
	if rawPath != "" {
		if strings.HasPrefix(rawPath, "http://") || strings.HasPrefix(rawPath, "https://") {
			return rawPath
		}
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "https://www.daemontalk.com"
		}
		baseURL = strings.TrimSuffix(baseURL, "/")
		if !strings.HasPrefix(rawPath, "/") {
			rawPath = "/" + rawPath
		}
		return fmt.Sprintf("%s%s", baseURL, rawPath)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://www.daemontalk.com"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s/blog/%s", baseURL, fallbackSlug)
}

// OSC52Copy generates the ANSI sequence to copy text into the SSH client's local clipboard
func OSC52Copy(text string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	return fmt.Sprintf("\x1b]52;c;%s\x07", encoded)
}

// OSC8Link creates a clickable hyperlink for modern terminal emulators
func OSC8Link(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

// openInBrowser opens a specified URL or file in the default system browser (for local desktop TUI)
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

package tui

import (
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

func OSC8Link(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

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

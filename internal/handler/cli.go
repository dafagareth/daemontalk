package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ANSI terminal color codes
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiGreen   = "\033[32m"
	ansiCyan    = "\033[36m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiWhite   = "\033[37m"
	ansiBgGray  = "\033[48;5;236m"
)

// IsCLIRequest checks whether the HTTP client is curl, wget, or requesting plain text.
func IsCLIRequest(r *http.Request) bool {
	if r.URL.Query().Get("format") == "text" || r.URL.Query().Get("plain") == "1" {
		return true
	}
	ua := strings.ToLower(r.UserAgent())
	if strings.Contains(ua, "curl") || strings.Contains(ua, "wget") || strings.Contains(ua, "httpie") {
		return true
	}
	if r.Header.Get("Accept") == "text/plain" {
		return true
	}
	return false
}

// stripANSI removes color codes if user requests nocolor
func applyColors(s string, enable bool) string {
	if enable {
		return s
	}
	replacements := []string{
		ansiReset, "", ansiBold, "", ansiDim, "", ansiGreen, "",
		ansiCyan, "", ansiYellow, "", ansiBlue, "", ansiMagenta, "",
		ansiWhite, "", ansiBgGray, "",
	}
	r := strings.NewReplacer(replacements...)
	return r.Replace(s)
}

// CLIMain handles "curl daemontalk.com" requests.
func (h *Handler) CLIMain(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("nocolor") != "1"
	posts := h.AllPosts()

	var b strings.Builder

	banner := `  ____                                _        _ _    
 |  _ \  __ _  ___ _ __ ___   ___  _ __ | |_ __ _| | | __
 | | | |/ _` + "`" + ` |/ _ \ '_ ` + "`" + ` _ \ / _ \| '_ \| __/ _` + "`" + ` | | |/ /
 | |_| | (_| |  __/ | | | | | (_) | | | | || (_| | |   < 
 |____/ \__,_|\___|_| |_| |_|\___/|_| |_|\__\__,_|_|_|\_\`

	b.WriteString(fmt.Sprintf("%s%s%s\n\n", ansiGreen+ansiBold, banner, ansiReset))
	b.WriteString(fmt.Sprintf(" %sdaemontalk · Daily Engineering & Linux Systems Dispatch%s\n", ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf(" %sHost:%s %s | %sTime:%s %s\n", ansiDim, ansiReset, r.Host, ansiDim, ansiReset, time.Now().UTC().Format("2006-01-02 15:04:05 UTC")))
	b.WriteString(fmt.Sprintf(" %s══════════════════════════════════════════════════════════════════════════%s\n\n", ansiDim, ansiReset))

	// Section 1: Latest Stories
	b.WriteString(fmt.Sprintf(" %s%s[ TOP DISPATCHES ]%s\n\n", ansiCyan, ansiBold, ansiReset))

	count := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		tags := strings.Join(p.Tags, ", ")
		b.WriteString(fmt.Sprintf("  %s•%s %s%s%s\n", ansiGreen, ansiReset, ansiBold, p.Title, ansiReset))
		b.WriteString(fmt.Sprintf("    %sDate:%s %s | %sTags:%s [%s]\n", ansiDim, ansiReset, p.Date.Format("2006-01-02"), ansiDim, ansiReset, tags))
		if p.Description != "" {
			b.WriteString(fmt.Sprintf("    %s%s%s\n", ansiDim, p.Description, ansiReset))
		}
		b.WriteString(fmt.Sprintf("    %sRead:%s curl -sL %s/p/%s\n\n", ansiYellow, ansiReset, r.Host, p.Slug))
		count++
		if count >= 6 {
			break
		}
	}

	// Section 3: CLI Shortcuts & API Guide
	b.WriteString(fmt.Sprintf(" %s%s[ TERMINAL ENDPOINTS ]%s\n\n", ansiMagenta, ansiBold, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/daily          %s# Stream full daily tech briefing%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/linux        %s# Stream Linux kernel & systems stories%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/go           %s# Stream Go concurrency & performance%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/python       %s# Stream Python & backend systems%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/recipes        %s# Production-ready Linux recipes%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/p/<slug>       %s# Read full article directly in terminal%s\n\n", r.Host, ansiDim, ansiReset))

	b.WriteString(fmt.Sprintf(" %sTip: Add '?nocolor=1' to disable ANSI color formatting.%s\n\n", ansiDim, ansiReset))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

// CLIDaily outputs a concise daily tech briefing.

// CLITag handles "curl daemontalk.com/t/{tag}".

// CLIPost handles "curl daemontalk.com/p/{slug}".

// CLIRecipes outputs curated eBPF, sysadmin & Linux troubleshooting one-liners.

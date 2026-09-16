package admin

import (
	"fmt"
	"strings"
	"time"

	"daemontalk/internal/post"
	"daemontalk/web/templates/shared"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fmtDate(p post.Post, _ ...string) string {
	d := p.Date
	return fmt.Sprintf("%s WIB", d.Format("02 January 2006"))
}

func fmtCommentTime(t time.Time, _ ...string) string {
	dur := time.Since(t)
	if dur < 0 {
		dur = 0
	}
	mins := int(dur.Minutes())
	if mins < 1 {
		return "just now"
	}
	if mins < 60 {
		return fmt.Sprintf("%dm ago", mins)
	}
	hours := int(dur.Hours())
	if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}
	days := hours / 24
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}
	return t.Format("02 Jan 2006")
}

func generateDigestText(posts []post.Post) string {
	var sb strings.Builder
	sb.WriteString("# ⚡ DaemonTalk Weekly Tech Digest\n\n")
	sb.WriteString("> Curated deep-dives in modern computing, AI, and universal technology.\n\n")
	sb.WriteString("## 📚 Featured Tech Dispatches This Week\n\n")

	base := strings.TrimSuffix(shared.SiteBaseURL, "/")
	count := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		sb.WriteString("### [" + p.Title + "](" + base + "/blog/" + p.Slug + ")\n")
		sb.WriteString(fmt.Sprintf("**Date:** %s  •  **Read Time:** %d min read\n\n", p.Date.Format("02 Jan 2006"), p.ReadTime))
		if p.Description != "" {
			sb.WriteString(p.Description + "\n\n")
		}
		sb.WriteString("[Read full dispatch →](" + base + "/blog/" + p.Slug + ")\n\n---\n\n")
		count++
		if count >= 5 {
			break
		}
	}

	sb.WriteString("## 📡 Syndicate & Feed\n\n")
	sb.WriteString("Subscribe to DaemonTalk dispatches via RSS or JSON Feed:\n")
	sb.WriteString("- [RSS Feed](" + base + "/rss.xml)\n")
	sb.WriteString("- [JSON Feed](" + base + "/feed.json)\n\n")
	sb.WriteString("---\n*Published by [DaemonTalk](" + base + ") · Independent Systems Journalism.*")

	return sb.String()
}

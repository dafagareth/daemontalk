package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"daemontalk/internal/post"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	renderCache sync.Map
)

func getContentPostsDir() string {
	candidates := []string{
		os.Getenv("CONTENT_DIR"),
		"content",
		"../content",
		"../../content",
	}

	for _, c := range candidates {
		if c == "" {
			continue
		}
		postsPath := filepath.Join(c, "posts")
		if info, err := os.Stat(postsPath); err == nil && info.IsDir() {
			return postsPath
		}
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return filepath.Join("content", "posts")
}

// FindPostFile locates the exact markdown file corresponding to a given post
func FindPostFile(p post.Post) string {
	postsDir := getContentPostsDir()

	// 1. Fast path: direct slug match
	directPath := filepath.Join(postsDir, p.Slug+".md")
	if _, err := os.Stat(directPath); err == nil {
		return directPath
	}

	// 2. Alias match
	for _, a := range p.Aliases {
		aliasPath := filepath.Join(postsDir, a+".md")
		if _, err := os.Stat(aliasPath); err == nil {
			return aliasPath
		}
	}

	// 3. Fallback: scan files only if direct paths were not found
	files, err := os.ReadDir(postsDir)
	if err != nil {
		return ""
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		path := filepath.Join(postsDir, f.Name())
		b, err := os.ReadFile(path)
		if err == nil && strings.Contains(string(b), p.Title) {
			return path
		}
	}

	return ""
}

// RenderPostMarkdown formats and renders the full markdown body using Glamour and theme styling
func RenderPostMarkdown(p post.Post, wrapWidth int, theme Theme) string {
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	cacheKey := fmt.Sprintf("%s:%d:%s", p.Slug, wrapWidth, theme.Name)
	if cached, ok := renderCache.Load(cacheKey); ok {
		if s, isStr := cached.(string); isStr && s != "" {
			return s
		}
	}

	exactFile := FindPostFile(p)

	var rawMD []byte
	if exactFile != "" {
		rawMD, _ = os.ReadFile(exactFile)
	}

	var content string
	if exactFile != "" && len(rawMD) > 0 {
		content = string(rawMD)
		// Remove frontmatter
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) == 3 {
				content = parts[2]
			}
		}
	} else if p.Description != "" {
		content = p.Description
	} else {
		content = "No content available for this dispatch."
	}

	// Format custom callouts into clean markdown quotes
	content = reCalloutOpen.ReplaceAllString(content, "\n> **[$1]** ")
	content = strings.ReplaceAll(content, "</callout>", "\n")

	// Format inline markdown images cleanly
	content = reMarkdownImage.ReplaceAllStringFunc(content, func(match string) string {
		sub := reMarkdownImage.FindStringSubmatch(match)
		if len(sub) == 3 {
			alt := sub[1]
			src := sub[2]
			if alt == "" {
				alt = "Illustration / Diagram"
			}
			return fmt.Sprintf("\n> 🖼️ **Image:** %s (`%s`)  •  *Press **'o'** to open in browser*\n", alt, src)
		}
		return match
	})

	// Graceful degradation for LaTeX Math blocks in Terminal
	content = reBlockMath.ReplaceAllString(content, "```math\n$1\n```")
	content = reInlineMath.ReplaceAllString(content, "`$1`")

	// Strip remaining raw HTML containers (e.g. <figure>, <details>, <summary>)
	content = reHTMLTags.ReplaceAllString(content, "")

	// 1. Glamour rendering for body
	glamourStyle := theme.GlamourStyle
	if glamourStyle == "" {
		glamourStyle = "dark"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(glamourStyle),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		r, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(wrapWidth),
		)
	}

	renderedBody := content
	if r != nil {
		if out, err := r.Render(content); err == nil {
			renderedBody = out
		}
	}

	// 2. Theme-Aware Header
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.TextNormal).
		Background(theme.SelectedBg).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.ActiveAccent)

	metaLabelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.TextNormal)

	metaBadgeStyle := lipgloss.NewStyle().
		Foreground(theme.ActiveAccent).
		Background(theme.SelectedBg).
		Padding(0, 1)

	dimStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted)

	dividerStyle := lipgloss.NewStyle().
		Foreground(theme.BorderInactive)

	var header strings.Builder
	header.WriteString("\n" + titleStyle.Render(p.Title) + "\n\n")

	dateBadge := metaBadgeStyle.Render(p.Date.Format("02 Jan 2006"))
	readBadge := metaBadgeStyle.Render(fmt.Sprintf("%d min", p.ReadTime))
	metaRow := fmt.Sprintf(" %s %s   •   %s %s",
		metaLabelStyle.Render("Date:"), dateBadge,
		metaLabelStyle.Render("Read Time:"), readBadge)

	if p.Author != "" {
		metaRow += fmt.Sprintf("   •   %s %s",
			metaLabelStyle.Render("Author:"), metaBadgeStyle.Render(p.Author))
	}
	header.WriteString(metaRow + "\n\n")

	if p.Cover != "" {
		coverBadge := metaBadgeStyle.Render(p.Cover)
		header.WriteString(fmt.Sprintf(" %s 🖼️  %s  %s\n\n",
			metaLabelStyle.Render("Cover Image:"),
			coverBadge,
			dimStyle.Render("(Press 'o' to open, 'w' for web)")))
	}

	if len(p.Tags) > 0 {
		var tagBadges []string
		for _, t := range p.Tags {
			tagBadges = append(tagBadges, metaBadgeStyle.Render(t))
		}
		header.WriteString(fmt.Sprintf(" %s  %s\n\n",
			metaLabelStyle.Render("Tags:"),
			strings.Join(tagBadges, " ")))
	}

	divWidth := wrapWidth
	if divWidth > 60 {
		divWidth = 60
	}
	header.WriteString(dividerStyle.Render(strings.Repeat("─", divWidth)) + "\n\n")

	fullOutput := header.String() + renderedBody
	renderCache.Store(cacheKey, fullOutput)
	return fullOutput
}

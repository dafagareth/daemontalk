package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"daemontalk/internal/post"
)

// FindPostFile locates the exact markdown file corresponding to a given post
func FindPostFile(p post.Post) string {
	files, err := os.ReadDir("content/posts")
	if err != nil {
		return ""
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		baseName := strings.TrimSuffix(f.Name(), ".md")

		// 1. Direct slug match
		if baseName == p.Slug {
			return filepath.Join("content/posts", f.Name())
		}

		// 2. Alias match
		for _, a := range p.Aliases {
			if baseName == a {
				return filepath.Join("content/posts", f.Name())
			}
		}

		// 3. Fallback: Title matching
		b, err := os.ReadFile(filepath.Join("content/posts", f.Name()))
		if err == nil && strings.Contains(string(b), p.Title) {
			return filepath.Join("content/posts", f.Name())
		}
	}

	return ""
}

// RenderPostMarkdown formats and renders the full markdown body using Glamour and theme styling
func RenderPostMarkdown(p post.Post, wrapWidth int, theme Theme) string {
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
	} else {
		content = p.Description
	}

	// Format inline markdown images cleanly without raw block art
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

	if wrapWidth < 20 {
		wrapWidth = 20
	}

	// 1. Glamour rendering for body
	glamourStyle := theme.GlamourStyle
	if glamourStyle == "" || glamourStyle == "tokyo-night" {
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

	// 2. Beautiful Theme-Aware Dynamic Header
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

	return header.String() + renderedBody
}

package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
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

// RenderPostMarkdown formats and renders the full markdown body using Glamour
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

	// Build clean article header
	var headerBuilder strings.Builder
	headerBuilder.WriteString(fmt.Sprintf("# %s\n\n", p.Title))
	headerBuilder.WriteString(fmt.Sprintf("**Date:** `%s`  •  **Read Time:** `%d min`", p.Date.Format("02 Jan 2006"), p.ReadTime))
	if p.Author != "" {
		headerBuilder.WriteString(fmt.Sprintf("  •  **Author:** `%s`", p.Author))
	}
	headerBuilder.WriteString("\n\n")

	if p.Cover != "" {
		headerBuilder.WriteString(fmt.Sprintf("**Cover Image:** 🖼️ `%s`  *(Press 'o' to open image, 'w' for web article)*\n\n", p.Cover))
	}

	if len(p.Tags) > 0 {
		tagBadges := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			tagBadges[i] = "`" + t + "`"
		}
		headerBuilder.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(tagBadges, " ")))
	}
	headerBuilder.WriteString("---\n\n")
	headerBuilder.WriteString(content)

	fullMarkdown := headerBuilder.String()

	if wrapWidth < 20 {
		wrapWidth = 20
	}

	// Glamour rendering with the active theme style
	glamourStyle := theme.GlamourStyle
	if glamourStyle == "" {
		glamourStyle = "dark"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(glamourStyle),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		// Fallback to standard dark if custom style isn't found
		r, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(wrapWidth),
		)
		if err != nil {
			return fullMarkdown
		}
	}

	out, err := r.Render(fullMarkdown)
	if err != nil {
		return fullMarkdown
	}

	return out
}

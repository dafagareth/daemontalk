package templates

import (
	"html"
	"html/template"
	"regexp"
	"sort"
	"strings"

	"daemontalk/internal/post"
)

var (
	cmBoldRe = regexp.MustCompile(`\*\*([^*\n]+)\*\*`)
	cmCodeRe = regexp.MustCompile("`([^`\n]+)`")
)

// renderCommentBody applies light-touch markdown (bold, code, newlines)
// to user-supplied comment text. HTML is escaped first for safety.
func renderCommentBody(body string) template.HTML {
	escaped := html.EscapeString(body)
	escaped = cmBoldRe.ReplaceAllString(escaped, "<strong>$1</strong>")
	escaped = cmCodeRe.ReplaceAllString(escaped, `<code class="font-mono text-xs bg-[var(--c-chip)] px-1 rounded">$1</code>`)
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")
	return template.HTML(escaped)
}

// highlightHTML HTML-escapes text then wraps all occurrences of query
// (case-insensitive) in <mark> tags for search result display.
func highlightHTML(text, query string) template.HTML {
	if query == "" || text == "" {
		return template.HTML(html.EscapeString(text))
	}
	lowerQ := strings.ToLower(query)
	lower := strings.ToLower(text)
	var buf strings.Builder
	pos := 0
	for {
		idx := strings.Index(lower[pos:], lowerQ)
		if idx < 0 {
			buf.WriteString(html.EscapeString(text[pos:]))
			break
		}
		buf.WriteString(html.EscapeString(text[pos : pos+idx]))
		buf.WriteString(`<mark class="bg-yellow-100 dark:bg-yellow-900/40 text-[var(--c-text)] rounded px-0.5">`)
		buf.WriteString(html.EscapeString(text[pos+idx : pos+idx+len(lowerQ)]))
		buf.WriteString("</mark>")
		pos += idx + len(lowerQ)
		if pos >= len(text) {
			break
		}
	}
	return template.HTML(buf.String())
}

type CategorySection struct {
	Tag      string
	LeadPost post.Post
	SubPosts []post.Post
}


// mapTagToCategory maps a post tag to the high-level main categories.
func mapTagToCategory(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	switch t {
	case "linux", "kernel", "ebpf":
		return "linux"
	case "systems", "rust", "go-compiler", "zig", "assembly", "compiler", "architecture":
		return "systems"
	case "backend", "go", "concurrency", "networking", "api", "web", "http", "caddy":
		return "backend"
	case "data", "python", "ai", "database", "storage", "clickhouse", "sqlite", "duckdb", "parquet":
		return "data"
	default:
		return t
	}
}

// getPortalCategories handles grouping of categories uniquely for the home page.
func getPortalCategories(posts []post.Post) ([]CategorySection, []CategorySection, []CategorySection) {
	seenSlugs := make(map[string]bool)

	// 1. Main Categories (top grid next to sidebar: linux, systems)
	mainTags := []string{"linux", "systems"}
	var mainCats []CategorySection
	for _, tag := range mainTags {
		matched := getPostsForCategory(posts, tag, seenSlugs)
		if len(matched) > 0 {
			mainCats = append(mainCats, buildSection(tag, matched))
		}
	}

	// 2. Middle Categories (devops, backend, data)
	midTags := []string{"devops", "backend", "data"}
	var midCats []CategorySection
	for _, tag := range midTags {
		matched := getPostsForCategory(posts, tag, seenSlugs)
		if len(matched) > 0 {
			midCats = append(midCats, buildSection(tag, matched))
		}
	}

	// 3. Other Categories (tools, performance, security, docker)
	otherTags := []string{"tools", "performance", "security", "docker"}
	var otherCats []CategorySection
	otherSeen := make(map[string]bool)
	for _, tag := range otherTags {
		matched := getPostsForCategory(posts, tag, otherSeen)
		if len(matched) > 0 {
			otherCats = append(otherCats, buildSection(tag, matched))
		}
	}

	return mainCats, midCats, otherCats
}

func getPostsForCategory(posts []post.Post, tag string, seenSlugs map[string]bool) []post.Post {
	var matched []post.Post
	for _, p := range posts {
		if seenSlugs[p.Slug] {
			continue
		}
		isMatch := false
		for _, t := range p.Tags {
			if mapTagToCategory(t) == tag {
				isMatch = true
				break
			}
		}
		if isMatch {
			matched = append(matched, p)
			seenSlugs[p.Slug] = true
		}
	}
	return matched
}

func buildSection(tag string, matched []post.Post) CategorySection {
	sec := CategorySection{
		Tag:      tag,
		LeadPost: matched[0],
	}
	if len(matched) > 1 {
		limit := 3
		if len(matched)-1 < limit {
			limit = len(matched) - 1
		}
		sec.SubPosts = matched[1 : limit+1]
	}
	return sec
}

func categoryTitle(cat CategorySection, lang string) string {
	switch strings.ToLower(cat.Tag) {
	case "systems":
		if lang == "id" {
			return "Rekayasa Sistem"
		}
		return "Systems Engineering"
	case "backend":
		if lang == "id" {
			return "Backend & Arsitektur"
		}
		return "Backend & Architecture"
	case "data":
		if lang == "id" {
			return "Data & AI"
		}
		return "Data & AI"
	case "linux":
		return "Linux"
	case "devops":
		return "DevOps"
	default:
		return strings.Title(cat.Tag)
	}
}

func getPopularPosts(posts []post.Post, viewCounts map[string]int, limit int) []post.Post {
	type ppost struct {
		p post.Post
		v int
	}
	var pp []ppost
	for _, p := range posts {
		pp = append(pp, ppost{p: p, v: viewCounts[p.Slug]})
	}
	sort.Slice(pp, func(i, j int) bool {
		return pp[i].v > pp[j].v
	})
	var out []post.Post
	for i := 0; i < len(pp) && i < limit; i++ {
		out = append(out, pp[i].p)
	}
	return out
}


type TagCount struct {
	Tag   string
	Count int
}

func getPopularTags(tagCounts map[string]int, limit int) []string {
	var tc []TagCount
	for tag, count := range tagCounts {
		tc = append(tc, TagCount{Tag: tag, Count: count})
	}
	sort.Slice(tc, func(i, j int) bool {
		if tc[i].Count == tc[j].Count {
			return tc[i].Tag < tc[j].Tag
		}
		return tc[i].Count > tc[j].Count
	})
	var out []string
	for i := 0; i < len(tc) && i < limit; i++ {
		out = append(out, tc[i].Tag)
	}
	return out
}

// generateDigestText creates a formatted Markdown weekly newsletter from published posts
func generateDigestText(posts []post.Post) string {
	var sb strings.Builder
	sb.WriteString("# ⚡ DaemonTalk Weekly Systems Digest\n\n")
	sb.WriteString("> Curated deep-dives in Linux Kernel architectures, language runtimes, and storage systems.\n\n")
	sb.WriteString("## 📚 Featured Engineering Dispatches This Week\n\n")

	count := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		sb.WriteString("### [" + p.Title + "](https://daemontalk.com/blog/" + p.Slug + ")\n")
		sb.WriteString("**Date:** " + p.Date.Format("02 Jan 2006") + "  •  **Read Time:** " + strings.TrimSpace(p.Description) + "\n\n")
		if p.Description != "" {
			sb.WriteString(p.Description + "\n\n")
		}
		sb.WriteString("[Read full dispatch →](https://daemontalk.com/blog/" + p.Slug + ")\n\n---\n\n")
		count++
		if count >= 5 {
			break
		}
	}

	sb.WriteString("## 💻 Read Terminal-First\n\n")
	sb.WriteString("Access DaemonTalk directly in your terminal:\n")
	sb.WriteString("```bash\n# Instant TUI over SSH\nssh daemontalk.com -p 2222\n\n# Daily dispatch stream\ncurl -sL daemontalk.com/daily\n```\n\n")
	sb.WriteString("---\n*Published by [DaemonTalk](https://daemontalk.com) · Independent Systems Journalism.*")

	return sb.String()
}

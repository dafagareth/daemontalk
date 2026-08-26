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
	if query == "" || text == "" || len(query) < 2 {
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

// getCategoryCluster finds posts matching a set of tags and builds a CategorySection.
// excludeSlugs allows skipping posts already displayed elsewhere on the homepage.
// By passing a pointer, the function automatically appends selected posts to prevent duplicates in subsequent calls.
func getCategoryCluster(posts []post.Post, categoryKey string, matchTags []string, subLimit int, excludeSlugs *[]string) CategorySection {
	tagSet := make(map[string]bool)
	for _, t := range matchTags {
		tagSet[strings.ToLower(strings.TrimSpace(t))] = true
	}

	excludeMap := make(map[string]bool)
	if excludeSlugs != nil {
		for _, s := range *excludeSlugs {
			excludeMap[s] = true
		}
	}

	var matched []post.Post
	for _, p := range posts {
		if excludeMap[p.Slug] {
			continue
		}
		for _, t := range p.Tags {
			cleanTag := strings.ToLower(strings.TrimSpace(t))
			if tagSet[cleanTag] {
				matched = append(matched, p)
				break
			}
		}
	}

	if len(matched) == 0 {
		// Fallback to any post if no exact match found
		for _, p := range posts {
			if !excludeMap[p.Slug] {
				matched = append(matched, p)
				break
			}
		}
	}

	sec := CategorySection{
		Tag: categoryKey,
	}
	if len(matched) > 0 {
		sec.LeadPost = matched[0]
		if excludeSlugs != nil {
			*excludeSlugs = append(*excludeSlugs, sec.LeadPost.Slug)
		}
	}
	if len(matched) > 1 {
		limit := subLimit
		if len(matched)-1 < limit {
			limit = len(matched) - 1
		}
		sec.SubPosts = matched[1 : limit+1]
		if excludeSlugs != nil {
			for _, p := range sec.SubPosts {
				*excludeSlugs = append(*excludeSlugs, p.Slug)
			}
		}
	}
	return sec
}

// collectSlugs gathers all slugs from CategorySections (lead + sub posts).
func collectSlugs(sections ...CategorySection) []string {
	var slugs []string
	for _, s := range sections {
		if s.LeadPost.Slug != "" {
			slugs = append(slugs, s.LeadPost.Slug)
		}
		for _, p := range s.SubPosts {
			slugs = append(slugs, p.Slug)
		}
	}
	return slugs
}

// collectPostSlugs gathers slugs from a slice of posts.
func collectPostSlugs(posts []post.Post) []string {
	var slugs []string
	for _, p := range posts {
		slugs = append(slugs, p.Slug)
	}
	return slugs
}

// getWirePosts returns quick headline dispatches avoiding already highlighted slugs.
func getWirePosts(posts []post.Post, excludeSlugs *[]string, limit int) []post.Post {
	excludeMap := make(map[string]bool)
	if excludeSlugs != nil {
		for _, s := range *excludeSlugs {
			excludeMap[s] = true
		}
	}
	var out []post.Post
	for _, p := range posts {
		if !excludeMap[p.Slug] {
			out = append(out, p)
			if excludeSlugs != nil {
				*excludeSlugs = append(*excludeSlugs, p.Slug)
			}
			if len(out) >= limit {
				break
			}
		}
	}
	return out
}

func categoryTitle(cat CategorySection, lang string) string {
	return categoryDisplayTitle(cat.Tag, lang)
}

func categoryDisplayTitle(key string, lang string) string {
	switch strings.ToLower(key) {
	case "linux":
		if lang == "id" {
			return "Linux & Kernel"
		}
		return "Linux & Kernel"
	case "go", "backend":
		if lang == "id" {
			return "Go & Backend"
		}
		return "Go & Backend"
	case "devops":
		if lang == "id" {
			return "DevOps & Cloud"
		}
		return "DevOps & Cloud"
	case "security":
		if lang == "id" {
			return "Security & Zero-Trust"
		}
		return "Security & Zero-Trust"
	case "systems":
		if lang == "id" {
			return "Systems & Low-Level"
		}
		return "Systems & Low-Level"
	case "performance":
		if lang == "id" {
			return "Performance & Memory"
		}
		return "Performance & Memory"
	case "docker":
		if lang == "id" {
			return "Docker & Containers"
		}
		return "Docker & Containers"
	case "tools":
		if lang == "id" {
			return "Tools & CLI"
		}
		return "Tools & CLI"
	case "terminal":
		if lang == "id" {
			return "Terminal & Workflow"
		}
		return "Terminal & Workflow"
	case "database":
		if lang == "id" {
			return "Database & Storage"
		}
		return "Database & Storage"
	case "networking":
		if lang == "id" {
			return "Networking & WireGuard"
		}
		return "Networking & WireGuard"
	case "ai":
		if lang == "id" {
			return "AI & Local LLM"
		}
		return "AI & Local LLM"
	case "ebpf":
		if lang == "id" {
			return "eBPF & Observability"
		}
		return "eBPF & Observability"
	case "debugging":
		if lang == "id" {
			return "Debugging & Dmesg"
		}
		return "Debugging & Dmesg"
	case "shell":
		if lang == "id" {
			return "Shell & Automation"
		}
		return "Shell & Automation"
	case "architecture":
		if lang == "id" {
			return "Architecture & Design"
		}
		return "Architecture & Design"
	case "storage":
		if lang == "id" {
			return "Storage & File Systems"
		}
		return "Storage & File Systems"
	default:
		return strings.Title(key)
	}
}

func getPopularPosts(posts []post.Post, viewCounts map[string]int, limit int) []post.Post {
	return getPopularPostsExcluding(posts, viewCounts, nil, limit)
}

func getPopularPostsExcluding(posts []post.Post, viewCounts map[string]int, excludeSlugs []string, limit int) []post.Post {
	excludeMap := make(map[string]bool)
	for _, s := range excludeSlugs {
		excludeMap[s] = true
	}

	type ppost struct {
		p post.Post
		v int
	}
	var pp []ppost
	for _, p := range posts {
		if !excludeMap[p.Slug] {
			pp = append(pp, ppost{p: p, v: viewCounts[p.Slug]})
		}
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
	sb.WriteString("```bash\n# Instant TUI over SSH\nssh ssh.daemontalk.com -p 2222\n\n# Daily dispatch stream\ncurl -sL daemontalk.com/daily\n```\n\n")
	sb.WriteString("---\n*Published by [DaemonTalk](https://daemontalk.com) · Independent Systems Journalism.*")

	return sb.String()
}

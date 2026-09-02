package templates

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"sort"
	"strings"

	"daemontalk/internal/post"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	escaped := html.EscapeString(text)
	escapedQ := html.EscapeString(query)
	lowerEscaped := strings.ToLower(escaped)
	lowerQ := strings.ToLower(escapedQ)

	qLen := len(lowerQ)
	if qLen == 0 {
		return template.HTML(escaped)
	}

	idx := strings.Index(lowerEscaped, lowerQ)
	if idx == -1 {
		return template.HTML(escaped)
	}

	var b strings.Builder
	b.Grow(len(escaped) + 64)
	lastIdx := 0

	for idx != -1 {
		b.WriteString(escaped[lastIdx:idx])
		b.WriteString(`<mark class="bg-yellow-100 dark:bg-yellow-900/40 text-[var(--c-text)] rounded px-0.5">`)
		b.WriteString(escaped[idx : idx+qLen])
		b.WriteString(`</mark>`)
		lastIdx = idx + qLen
		nextMatch := strings.Index(lowerEscaped[lastIdx:], lowerQ)
		if nextMatch == -1 {
			break
		}
		idx = lastIdx + nextMatch
	}
	b.WriteString(escaped[lastIdx:])
	return template.HTML(b.String())
}

type CategorySection struct {
	Tag      string
	LeadPost post.Post
	SubPosts []post.Post
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

type CategoryTagLink struct {
	Name string
	Slug string
}

func getCategoryTagLinks(key string) []CategoryTagLink {
	switch strings.ToLower(key) {
	case "software":
		return []CategoryTagLink{
			{Name: "Software", Slug: "software"},
			{Name: "Development", Slug: "development"},
		}
	case "linux":
		return []CategoryTagLink{
			{Name: "Linux", Slug: "linux"},
			{Name: "OS", Slug: "os"},
		}
	case "ai":
		return []CategoryTagLink{
			{Name: "AI", Slug: "ai"},
			{Name: "Machine Learning", Slug: "machine-learning"},
		}
	case "security":
		return []CategoryTagLink{
			{Name: "Security", Slug: "security"},
			{Name: "Privacy", Slug: "privacy"},
		}
	case "networking":
		return []CategoryTagLink{
			{Name: "Networking", Slug: "networking"},
			{Name: "Protocols", Slug: "protocols"},
		}
	case "gaming":
		return []CategoryTagLink{
			{Name: "Gaming", Slug: "gaming"},
			{Name: "Graphics", Slug: "graphics"},
		}
	case "tools":
		return []CategoryTagLink{
			{Name: "Tools", Slug: "tools"},
			{Name: "Workflow", Slug: "workflow"},
		}
	case "science":
		return []CategoryTagLink{
			{Name: "Science", Slug: "science"},
			{Name: "Research", Slug: "research"},
		}
	case "policy":
		return []CategoryTagLink{
			{Name: "Tech Policy", Slug: "policy"},
			{Name: "Law", Slug: "law"},
		}
	case "backend-architecture", "backend":
		return []CategoryTagLink{
			{Name: "Backend", Slug: "backend"},
			{Name: "Architecture", Slug: "architecture"},
		}
	case "systems":
		return []CategoryTagLink{
			{Name: "Systems", Slug: "systems"},
			{Name: "Low-Level", Slug: "low-level"},
		}
	case "container-internals", "containers":
		return []CategoryTagLink{
			{Name: "Containers", Slug: "containers"},
			{Name: "Internals", Slug: "container-internals"},
		}
	case "terminal":
		return []CategoryTagLink{
			{Name: "Terminal", Slug: "terminal"},
			{Name: "Shell", Slug: "shell"},
		}
	case "database":
		return []CategoryTagLink{
			{Name: "Database", Slug: "database"},
			{Name: "Storage", Slug: "storage"},
		}
	case "ebpf":
		return []CategoryTagLink{
			{Name: "eBPF", Slug: "ebpf"},
			{Name: "Observability", Slug: "observability"},
		}
	case "crypto":
		return []CategoryTagLink{
			{Name: "Privacy", Slug: "privacy"},
			{Name: "Crypto", Slug: "crypto"},
		}
	default:
		return []CategoryTagLink{
			{Name: cases.Title(language.English).String(key), Slug: key},
		}
	}
}

func categoryTitle(cat CategorySection, lang string) string {
	return categoryDisplayTitle(cat.Tag, lang)
}

func categoryDisplayTitle(key string, lang string) string {
	switch strings.ToLower(key) {
	case "software":
		return "Software & Development"
	case "linux":
		return "Linux & OS"
	case "ai":
		return "AI & Machine Learning"
	case "security":
		return "Security & Privacy"
	case "networking":
		return "Networking & Protocols"
	case "gaming":
		return "Gaming & Graphics"
	case "tools":
		return "Tools & Workflow"
	case "science":
		return "Science & Research"
	case "policy":
		return "Tech Policy & Law"
	case "backend-architecture":
		return "Backend Architecture"
	case "systems":
		return "Systems & Low-Level"
	case "container-internals":
		return "Container Internals"
	case "terminal":
		return "Terminal & Shell"
	case "database":
		return "Database & Storage"
	case "ebpf":
		return "eBPF & Observability"
	case "crypto":
		return "Privacy & Crypto"
	default:
		return cases.Title(language.English).String(key)
	}
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
	sb.WriteString("# ⚡ DaemonTalk Weekly Tech Digest\n\n")
	sb.WriteString("> Curated deep-dives in modern computing, AI, and universal technology.\n\n")
	sb.WriteString("## 📚 Featured Tech Dispatches This Week\n\n")

	base := strings.TrimSuffix(SiteBaseURL, "/")
	hostDomain := strings.TrimPrefix(strings.TrimPrefix(base, "https://"), "http://")
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

	sb.WriteString("## 💻 Read Terminal-First\n\n")
	sb.WriteString("Access DaemonTalk directly in your terminal:\n")
	sb.WriteString("```bash\n# Instant TUI over SSH\nssh ssh.daemontalk.com -p 2222\n\n# Daily dispatch stream\ncurl -sL " + hostDomain + "/daily\n```\n\n")
	sb.WriteString("---\n*Published by [DaemonTalk](" + base + ") · Independent Systems Journalism.*")

	return sb.String()
}

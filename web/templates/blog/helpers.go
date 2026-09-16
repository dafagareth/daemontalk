package blog

import (
	"html"
	"html/template"
	"regexp"
	"sort"
	"strings"

	"daemontalk/internal/post"
)

type tagStat struct {
	Tag   string
	Count int
}

type PostNav struct {
	HasPrev bool
	HasNext bool
	Prev    post.Post
	Next    post.Post
}

func postDifficulty(p post.Post) string {
	for _, t := range p.Tags {
		tl := strings.ToLower(t)
		if tl == "architecture" || tl == "debugging" || tl == "storage" || tl == "cgroups" || tl == "profiling" || tl == "concurrency" || tl == "performance" {
			return "DEEP DIVE"
		}
	}
	if p.ReadTime >= 5 {
		return "DEEP DIVE"
	}
	if p.ReadTime >= 3 {
		return "INTERMEDIATE"
	}
	return "BEGINNER"
}

func sortedTags(m map[string]int) []tagStat {
	out := make([]tagStat, 0, len(m))
	for t, c := range m {
		out = append(out, tagStat{Tag: t, Count: c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Tag < out[j].Tag
	})
	return out
}

func searchIndex(p post.Post) string {
	return strings.ToLower(p.Title + " " + strings.Join(p.Tags, " ") + " " + p.Description)
}

func postCover(p post.Post) string {
	return p.Cover
}

var (
	cmBoldRe = regexp.MustCompile(`\*\*([^*\n]+)\*\*`)
	cmCodeRe = regexp.MustCompile("`([^`\n]+)`")
)

func renderCommentBody(body string) template.HTML {
	escaped := html.EscapeString(body)
	escaped = cmBoldRe.ReplaceAllString(escaped, "<strong>$1</strong>")
	escaped = cmCodeRe.ReplaceAllString(escaped, `<code class="font-mono text-xs bg-[var(--c-chip)] px-1 rounded">$1</code>`)
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")
	return template.HTML(escaped)
}

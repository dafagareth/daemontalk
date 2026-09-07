package templates

import (
	"sort"
	"strings"

	"daemontalk/internal/post"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CategorySection struct {
	Tag      string
	LeadPost post.Post
	SubPosts []post.Post
}

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

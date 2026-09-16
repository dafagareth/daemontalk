package portal

import (
	"sort"
	"strings"

	"daemontalk/internal/post"
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

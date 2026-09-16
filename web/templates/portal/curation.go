package portal

import (
	"sort"
	"strings"
	"time"

	"daemontalk/internal/post"
)

func getSupportingPosts(posts []post.Post, lead post.Post, excludeSlugs *[]string, limit int) []post.Post {
	excludeMap := make(map[string]bool)
	if excludeSlugs != nil {
		for _, s := range *excludeSlugs {
			excludeMap[s] = true
		}
	}

	tagSet := make(map[string]bool)
	for _, t := range lead.Tags {
		clean := strings.ToLower(strings.TrimSpace(t))
		if clean != "" {
			tagSet[clean] = true
		}
	}

	var out []post.Post
	if len(tagSet) > 0 {
		for _, p := range posts {
			if excludeMap[p.Slug] {
				continue
			}
			for _, t := range p.Tags {
				if tagSet[strings.ToLower(strings.TrimSpace(t))] {
					out = append(out, p)
					excludeMap[p.Slug] = true
					if excludeSlugs != nil {
						*excludeSlugs = append(*excludeSlugs, p.Slug)
					}
					break
				}
			}
			if len(out) >= limit {
				return out
			}
		}
	}

	for _, p := range posts {
		if len(out) >= limit {
			break
		}
		if !excludeMap[p.Slug] {
			out = append(out, p)
			excludeMap[p.Slug] = true
			if excludeSlugs != nil {
				*excludeSlugs = append(*excludeSlugs, p.Slug)
			}
		}
	}
	return out
}

func getPopularPostsThisWeek(posts []post.Post, viewCounts map[string]int, excludeSlugs *[]string, limit int) []post.Post {
	excludeMap := make(map[string]bool)
	if excludeSlugs != nil {
		for _, s := range *excludeSlugs {
			excludeMap[s] = true
		}
	}

	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	type ppost struct {
		p post.Post
		v int
	}
	var weekPosts []ppost
	var allCandidates []ppost

	for _, p := range posts {
		if excludeMap[p.Slug] {
			continue
		}
		item := ppost{p: p, v: viewCounts[p.Slug]}
		allCandidates = append(allCandidates, item)
		if p.Date.After(sevenDaysAgo) {
			weekPosts = append(weekPosts, item)
		}
	}

	sortPopular := func(slice []ppost) {
		sort.Slice(slice, func(i, j int) bool {
			if slice[i].v == slice[j].v {
				return slice[i].p.Date.After(slice[j].p.Date)
			}
			return slice[i].v > slice[j].v
		})
	}

	sortPopular(weekPosts)
	sortPopular(allCandidates)

	var out []post.Post
	for _, wp := range weekPosts {
		if len(out) >= limit {
			break
		}
		out = append(out, wp.p)
		excludeMap[wp.p.Slug] = true
		if excludeSlugs != nil {
			*excludeSlugs = append(*excludeSlugs, wp.p.Slug)
		}
	}

	for _, cp := range allCandidates {
		if len(out) >= limit {
			break
		}
		if !excludeMap[cp.p.Slug] {
			out = append(out, cp.p)
			excludeMap[cp.p.Slug] = true
			if excludeSlugs != nil {
				*excludeSlugs = append(*excludeSlugs, cp.p.Slug)
			}
		}
	}

	return out
}

func getPopularPosts(posts []post.Post, viewCounts map[string]int, excludeSlugs *[]string, limit int) []post.Post {
	excludeMap := make(map[string]bool)
	if excludeSlugs != nil {
		for _, s := range *excludeSlugs {
			excludeMap[s] = true
		}
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
		if pp[i].v == pp[j].v {
			return pp[i].p.Date.After(pp[j].p.Date)
		}
		return pp[i].v > pp[j].v
	})
	var out []post.Post
	for i := 0; i < len(pp) && i < limit; i++ {
		out = append(out, pp[i].p)
		if excludeSlugs != nil {
			*excludeSlugs = append(*excludeSlugs, pp[i].p.Slug)
		}
	}
	return out
}

func getRemainingPosts(posts []post.Post, excludeSlugs []string) []post.Post {
	excludeMap := make(map[string]bool, len(excludeSlugs))
	for _, s := range excludeSlugs {
		excludeMap[s] = true
	}
	var out []post.Post
	for _, p := range posts {
		if !excludeMap[p.Slug] {
			out = append(out, p)
		}
	}
	return out
}

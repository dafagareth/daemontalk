package graph

import (
	"encoding/json"
	"hash/fnv"
	"regexp"
	"sort"
	"strings"

	"daemontalk/internal/post"
)

type NodeType string

const (
	NodePost NodeType = "post"
	NodeTag  NodeType = "tag"
)

type Node struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Type        NodeType `json:"type"`
	Slug        string   `json:"slug,omitempty"`
	Tag         string   `json:"tag,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Date        string   `json:"date,omitempty"`
	ReadTime    int      `json:"readTime,omitempty"`
	Description string   `json:"description,omitempty"`
	Radius      float64  `json:"r"`
	Color       string   `json:"color"`
	LinkCount   int      `json:"linkCount"`
}

type Link struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
	Kind   string  `json:"kind"`
}

type GraphStats struct {
	TotalPosts int `json:"totalPosts"`
	TotalTags  int `json:"totalTags"`
	TotalLinks int `json:"totalLinks"`
}

type GraphData struct {
	Nodes []Node     `json:"nodes"`
	Links []Link     `json:"links"`
	Stats GraphStats `json:"stats"`
}

func (g GraphData) ToJSON() string {
	b, err := json.Marshal(g)
	if err != nil {
		return "{}"
	}
	return string(b)
}

var internalLinkRe = regexp.MustCompile(`(?:href="|\]\()(?:(?:https?://[^/]+)?(?:/id)?/blog/([a-zA-Z0-9_-]+))`)

var tagPalette = map[string]string{
	"linux":        "#38bdf8",
	"kernel":       "#0284c7",
	"ebpf":         "#06b6d4",
	"go":           "#00add8",
	"golang":       "#00add8",
	"rust":         "#f97316",
	"security":     "#ef4444",
	"storage":      "#10b981",
	"database":     "#059669",
	"cgroups":      "#8b5cf6",
	"docker":       "#a855f7",
	"containers":   "#a855f7",
	"ai":           "#eab308",
	"performance":  "#f59e0b",
	"concurrency":  "#6366f1",
	"networking":   "#14b8a6",
	"architecture": "#94a3b8",
	"tools":        "#64748b",
	"software":     "#3b82f6",
	"gaming":       "#ec4899",
	"science":      "#10b981",
}

var fallbackColors = []string{
	"#38bdf8", "#818cf8", "#a855f7", "#ec4899",
	"#f43f5e", "#f97316", "#eab308", "#10b981",
	"#14b8a6", "#06b6d4", "#64748b",
}

func tagColor(t string) string {
	low := strings.ToLower(strings.TrimSpace(t))
	if c, ok := tagPalette[low]; ok {
		return c
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(low))
	idx := int(h.Sum32()) % len(fallbackColors)
	if idx < 0 {
		idx = -idx
	}
	return fallbackColors[idx]
}

func Build(posts []post.Post) GraphData {
	postSlugMap := make(map[string]post.Post)
	tagCountMap := make(map[string]int)

	for _, p := range posts {
		if p.Draft {
			continue
		}
		postSlugMap[p.Slug] = p
		for _, t := range p.Tags {
			cleaned := strings.ToLower(strings.TrimSpace(t))
			if cleaned != "" {
				tagCountMap[cleaned]++
			}
		}
	}

	nodesMap := make(map[string]Node)
	var links []Link
	seenLinks := make(map[string]bool)

	var tagNames []string
	for t := range tagCountMap {
		tagNames = append(tagNames, t)
	}
	sort.Strings(tagNames)

	for _, t := range tagNames {
		cnt := tagCountMap[t]
		r := 6.0 + float64(cnt)*0.6
		if r > 10.0 {
			r = 10.0
		}
		nodeID := "tag:" + t
		nodesMap[nodeID] = Node{
			ID:     nodeID,
			Label:  "#" + t,
			Type:   NodeTag,
			Tag:    t,
			Radius: r,
			Color:  tagColor(t),
		}
	}

	for _, p := range posts {
		if p.Draft {
			continue
		}

		postNodeID := "post:" + p.Slug
		primaryTag := ""
		if len(p.Tags) > 0 {
			primaryTag = strings.ToLower(strings.TrimSpace(p.Tags[0]))
		}

		r := 4.5 + float64(p.ReadTime)*0.2
		if r > 7.0 {
			r = 7.0
		}

		color := tagColor(primaryTag)
		if primaryTag == "" {
			color = "#94a3b8"
		}

		dateStr := ""
		if !p.Date.IsZero() {
			dateStr = p.Date.Format("02 Jan 2006")
		}

		nodesMap[postNodeID] = Node{
			ID:          postNodeID,
			Label:       p.Title,
			Type:        NodePost,
			Slug:        p.Slug,
			Tag:         primaryTag,
			Tags:        p.Tags,
			Date:        dateStr,
			ReadTime:    p.ReadTime,
			Description: p.Description,
			Radius:      r,
			Color:       color,
		}

		for _, t := range p.Tags {
			cleanTag := strings.ToLower(strings.TrimSpace(t))
			tagNodeID := "tag:" + cleanTag
			if _, ok := nodesMap[tagNodeID]; ok {
				linkKey := postNodeID + "<->" + tagNodeID
				if !seenLinks[linkKey] {
					seenLinks[linkKey] = true
					links = append(links, Link{
						Source: postNodeID,
						Target: tagNodeID,
						Weight: 1.0,
						Kind:   "tag",
					})
				}
			}
		}

		bodyStr := string(p.Body)
		matches := internalLinkRe.FindAllStringSubmatch(bodyStr, -1)
		for _, m := range matches {
			if len(m) > 1 {
				targetSlug := m[1]
				if targetSlug != p.Slug {
					if _, exists := postSlugMap[targetSlug]; exists {
						targetNodeID := "post:" + targetSlug

						var k1, k2 string
						if postNodeID < targetNodeID {
							k1, k2 = postNodeID, targetNodeID
						} else {
							k1, k2 = targetNodeID, postNodeID
						}
						linkKey := k1 + "<=>" + k2
						if !seenLinks[linkKey] {
							seenLinks[linkKey] = true
							links = append(links, Link{
								Source: postNodeID,
								Target: targetNodeID,
								Weight: 2.2,
								Kind:   "crosslink",
							})
						}
					}
				}
			}
		}
	}

	for _, l := range links {
		if n, ok := nodesMap[l.Source]; ok {
			n.LinkCount++
			nodesMap[l.Source] = n
		}
		if n, ok := nodesMap[l.Target]; ok {
			n.LinkCount++
			nodesMap[l.Target] = n
		}
	}

	var nodes []Node
	for _, n := range nodesMap {
		nodes = append(nodes, n)
	}

	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type == NodeTag
		}
		return nodes[i].Label < nodes[j].Label
	})

	stats := GraphStats{
		TotalPosts: len(postSlugMap),
		TotalTags:  len(tagNames),
		TotalLinks: len(links),
	}

	return GraphData{
		Nodes: nodes,
		Links: links,
		Stats: stats,
	}
}

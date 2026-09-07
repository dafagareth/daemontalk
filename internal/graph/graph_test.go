package graph

import (
	"html/template"
	"testing"
	"time"

	"daemontalk/internal/post"
)

func TestBuildGraph(t *testing.T) {
	posts := []post.Post{
		{
			Title:       "Deep Dive into eBPF Tracing",
			Slug:        "ebpf-tracing",
			Tags:        []string{"linux", "ebpf"},
			Date:        time.Now(),
			ReadTime:    5,
			Description: "An in-depth look at Linux eBPF internals.",
			Body:        template.HTML(`<p>Read also about <a href="/blog/io-uring">io_uring</a>.</p>`),
		},
		{
			Title:       "High Performance I/O with io_uring",
			Slug:        "io-uring",
			Tags:        []string{"linux"},
			Date:        time.Now().Add(-time.Hour),
			ReadTime:    4,
			Description: "Modern Linux asynchronous I/O with io_uring.",
			Body:        template.HTML(`<p>See our post on [eBPF](/blog/ebpf-tracing).</p>`),
		},
		{
			Title:       "Go Memory Model and Concurrency",
			Slug:        "go-memory-model",
			Tags:        []string{"go", "concurrency"},
			Date:        time.Now().Add(-2 * time.Hour),
			ReadTime:    6,
			Description: "Understanding happens-before in Go.",
			Body:        template.HTML(`<p>Go runtime details.</p>`),
		},
	}

	data := Build(posts)

	if data.Stats.TotalPosts != 3 {
		t.Fatalf("expected 3 posts, got %d", data.Stats.TotalPosts)
	}
	if data.Stats.TotalTags != 4 {
		t.Fatalf("expected 4 tags, got %d", data.Stats.TotalTags)
	}

	postFound := false
	for _, n := range data.Nodes {
		if n.ID == "post:ebpf-tracing" {
			postFound = true
			if n.Type != NodePost {
				t.Errorf("expected node type post, got %s", n.Type)
			}
			if n.Tag != "linux" {
				t.Errorf("expected primary tag linux, got %s", n.Tag)
			}
			break
		}
	}
	if !postFound {
		t.Error("post:ebpf-tracing node not found")
	}

	crosslinkFound := false
	for _, l := range data.Links {
		if l.Kind == "crosslink" {
			if (l.Source == "post:ebpf-tracing" && l.Target == "post:io-uring") ||
				(l.Source == "post:io-uring" && l.Target == "post:ebpf-tracing") {
				crosslinkFound = true
				if l.Weight < 2.0 {
					t.Errorf("expected higher weight for crosslink, got %f", l.Weight)
				}
				break
			}
		}
	}
	if !crosslinkFound {
		t.Error("expected direct crosslink between ebpf-tracing and io-uring")
	}

	jsonStr := data.ToJSON()
	if jsonStr == "" || jsonStr == "{}" {
		t.Error("expected valid JSON string from ToJSON()")
	}
}

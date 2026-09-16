package portal

import (
	"testing"
	"time"

	"daemontalk/internal/post"
)

func TestGetSupportingPosts(t *testing.T) {
	now := time.Now()
	lead := post.Post{
		Slug: "lead-opinion",
		Tags: []string{"opinion", "career"},
		Date: now,
	}

	posts := []post.Post{
		lead,
		{Slug: "rust-post", Tags: []string{"rust"}, Date: now.Add(-1 * time.Hour)},
		{Slug: "opinion-sub1", Tags: []string{"opinion"}, Date: now.Add(-2 * time.Hour)},
		{Slug: "go-post", Tags: []string{"go"}, Date: now.Add(-3 * time.Hour)},
		{Slug: "career-sub2", Tags: []string{"career"}, Date: now.Add(-4 * time.Hour)},
	}

	var used []string
	used = append(used, lead.Slug)

	subs := getSupportingPosts(posts, lead, &used, 2)
	if len(subs) != 2 {
		t.Fatalf("expected 2 supporting posts, got %d", len(subs))
	}
	if subs[0].Slug != "opinion-sub1" || subs[1].Slug != "career-sub2" {
		t.Errorf("unexpected supporting posts: got [%s, %s], expected [opinion-sub1, career-sub2]", subs[0].Slug, subs[1].Slug)
	}
	if len(used) != 3 {
		t.Errorf("expected 3 used slugs, got %d", len(used))
	}
}

func TestGetSupportingPostsFallback(t *testing.T) {
	now := time.Now()
	lead := post.Post{
		Slug: "lead-unique",
		Tags: []string{"unique-tag"},
		Date: now,
	}

	posts := []post.Post{
		lead,
		{Slug: "post-1", Tags: []string{"rust"}, Date: now.Add(-1 * time.Hour)},
		{Slug: "post-2", Tags: []string{"go"}, Date: now.Add(-2 * time.Hour)},
	}

	var used []string
	used = append(used, lead.Slug)

	subs := getSupportingPosts(posts, lead, &used, 2)
	if len(subs) != 2 {
		t.Fatalf("expected 2 supporting posts fallback, got %d", len(subs))
	}
	if subs[0].Slug != "post-1" || subs[1].Slug != "post-2" {
		t.Errorf("unexpected fallback posts: [%s, %s]", subs[0].Slug, subs[1].Slug)
	}
}

func TestGetPopularPostsThisWeek(t *testing.T) {
	now := time.Now()
	posts := []post.Post{
		{Slug: "old-viral", Title: "Old Viral", Date: now.AddDate(0, 0, -30)},
		{Slug: "week-low", Title: "Week Low", Date: now.AddDate(0, 0, -2)},
		{Slug: "week-high", Title: "Week High", Date: now.AddDate(0, 0, -1)},
	}
	views := map[string]int{
		"old-viral": 9999,
		"week-low":  20,
		"week-high": 150,
	}

	var used []string
	pop := getPopularPostsThisWeek(posts, views, &used, 2)
	if len(pop) != 2 {
		t.Fatalf("expected 2 week popular posts, got %d", len(pop))
	}
	if pop[0].Slug != "week-high" || pop[1].Slug != "week-low" {
		t.Errorf("expected [week-high, week-low], got [%s, %s]", pop[0].Slug, pop[1].Slug)
	}
}

func TestGetPopularPosts(t *testing.T) {
	now := time.Now()
	posts := []post.Post{
		{Slug: "p1", Title: "P1", Date: now.Add(-1 * time.Hour)},
		{Slug: "p2", Title: "P2", Date: now},
		{Slug: "p3", Title: "P3", Date: now.Add(-2 * time.Hour)},
	}
	views := map[string]int{
		"p1": 50,
		"p2": 10,
		"p3": 100,
	}

	var used []string
	pop := getPopularPosts(posts, views, &used, 2)
	if len(pop) != 2 {
		t.Fatalf("expected 2 popular posts, got %d", len(pop))
	}
	if pop[0].Slug != "p3" || pop[1].Slug != "p1" {
		t.Errorf("expected [p3, p1], got [%s, %s]", pop[0].Slug, pop[1].Slug)
	}
}

func TestGetRemainingPosts(t *testing.T) {
	posts := []post.Post{
		{Slug: "p1"},
		{Slug: "p2"},
		{Slug: "p3"},
		{Slug: "p4"},
	}
	exclude := []string{"p1", "p3"}
	rem := getRemainingPosts(posts, exclude)
	if len(rem) != 2 {
		t.Fatalf("expected 2 remaining posts, got %d", len(rem))
	}
	if rem[0].Slug != "p2" || rem[1].Slug != "p4" {
		t.Errorf("expected [p2, p4], got [%s, %s]", rem[0].Slug, rem[1].Slug)
	}
}

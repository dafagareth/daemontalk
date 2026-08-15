package post

import (
	"strings"
	"testing"
	"time"
)

func TestLoadAll(t *testing.T) {
	posts, err := LoadAll("testdata")
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post (draft excluded), got %d", len(posts))
	}

	p := posts[0]
	if p.Title != "Hello World" {
		t.Errorf("title: got %q, want %q", p.Title, "Hello World")
	}
	if p.Slug != "hello-world" {
		t.Errorf("slug: got %q, want %q", p.Slug, "hello-world")
	}
	if p.Lang != "en" {
		t.Errorf("lang: got %q, want %q", p.Lang, "en")
	}
	if p.Draft {
		t.Error("draft should be false")
	}
	if p.Date != time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("date: got %v", p.Date)
	}
	if len(p.Tags) != 2 {
		t.Errorf("tags: expected 2, got %d", len(p.Tags))
	}
	if p.ReadTime < 1 {
		t.Errorf("read time should be at least 1, got %d", p.ReadTime)
	}
	if p.Body == "" {
		t.Error("body should not be empty")
	}
}

func TestLoadAll_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	posts, err := LoadAll(dir)
	if err != nil {
		t.Fatalf("LoadAll empty dir: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestLoadAll_NonExistentDir(t *testing.T) {
	posts, err := LoadAll("/no/such/dir")
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestFindBySlug(t *testing.T) {
	posts, _ := LoadAll("testdata")
	p, ok := FindBySlug(posts, "hello-world")
	if !ok {
		t.Fatal("expected to find hello-world")
	}
	if p.Title != "Hello World" {
		t.Errorf("wrong post returned: %q", p.Title)
	}

	_, ok = FindBySlug(posts, "nonexistent")
	if ok {
		t.Error("should not find nonexistent slug")
	}
}

func TestReadTime(t *testing.T) {
	short := readTime("<p>hello world</p>")
	if short != 1 {
		t.Errorf("short text should be 1 min, got %d", short)
	}

	long := ""
	for i := 0; i < 300; i++ {
		long += "word "
	}
	mins := readTime(long)
	if mins < 1 {
		t.Errorf("300 words should be at least 1 min, got %d", mins)
	}
}

func TestParseBytes(t *testing.T) {
	src := []byte(`---
title: "Post dari DB"
slug: post-dari-db
date: 2026-07-04
lang: id
draft: true
tags: ["go", "web"]
---

Paragraf pembuka untuk deskripsi.

## Bagian Pertama

Isi bagian pertama.
`)
	p, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if p.Title != "Post dari DB" || p.Slug != "post-dari-db" {
		t.Errorf("title/slug salah: %q %q", p.Title, p.Slug)
	}
	if !p.Draft {
		t.Error("draft harus true")
	}
	if len(p.Tags) != 2 {
		t.Errorf("tags: %v", p.Tags)
	}
	if len(p.TOC) != 1 || p.TOC[0].Title != "Bagian Pertama" {
		t.Errorf("TOC: %+v", p.TOC)
	}
	if p.ReadTime < 1 {
		t.Errorf("readtime: %d", p.ReadTime)
	}
	if !strings.Contains(string(p.Body), "<h2") {
		t.Error("body harus mengandung h2")
	}
}

func TestAliasesAndFindBySlug(t *testing.T) {
	src := []byte(`---
title: "Arch Linux Tips"
slug: 803461ff
aliases: ["archlinux-funfact-tips", "arch-tips"]
---
Hello Arch!
`)
	p, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(p.Aliases) != 2 || p.Aliases[0] != "archlinux-funfact-tips" {
		t.Errorf("aliases not parsed correctly: %v", p.Aliases)
	}

	posts := []Post{p}
	found, ok := FindBySlug(posts, "803461ff")
	if !ok || found.Slug != "803461ff" {
		t.Errorf("expected to find by canonical slug")
	}

	foundAlias, ok := FindBySlug(posts, "archlinux-funfact-tips")
	if !ok || foundAlias.Slug != "803461ff" {
		t.Errorf("expected to find by alias")
	}
}

func TestCarouselAndGalleryParsing(t *testing.T) {
	src := []byte(`---
title: "Testing Carousel & Gallery"
slug: test-carousel
date: 2026-08-14
---

Here is a carousel:

` + "```carousel" + `
![Diagram 1](/static/images/1.png "Caption 1")
![Diagram 2](/static/images/2.png "Caption 2")
` + "```" + `

Here is a gallery:

` + "```gallery" + `
![Before](/static/images/before.png "Before State")
![After](/static/images/after.png "After State")
` + "```" + `

Here is an FAQ:

` + "```faq" + `
Q: Mengapa memilih Linux epoll?
A: Karena skalabilitas O(1) dengan event-driven I/O.

Q: Apakah support multi-thread?
A: Ya, thread-safe dengan ` + "`sync.Mutex`" + `.
` + "```" + `

Here is an Author card:

` + "```author" + `
name: Dafa Gareth
role: Software Engineer
avatar: /static/logo/logo-dark.png
bio: Systems programmer specializing in Linux and Go.
github: https://github.com/dafagareth
email: dafagareth@gmail.com
` + "```" + `
`)

	p, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	bodyStr := string(p.Body)
	if !strings.Contains(bodyStr, "post-carousel-wrap") {
		t.Errorf("expected post-carousel-wrap class in body, got: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "1 / 2") || !strings.Contains(bodyStr, "2 / 2") {
		t.Errorf("expected carousel counter badges in body")
	}
	if !strings.Contains(bodyStr, "Caption 1") || !strings.Contains(bodyStr, "Caption 2") {
		t.Errorf("expected carousel captions in body")
	}
	if !strings.Contains(bodyStr, "post-gallery-wrap") {
		t.Errorf("expected post-gallery-wrap class in body, got: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "Before State") || !strings.Contains(bodyStr, "After State") {
		t.Errorf("expected gallery captions in body")
	}
	if !strings.Contains(bodyStr, "post-faq-wrap") || !strings.Contains(bodyStr, "Mengapa memilih Linux epoll?") {
		t.Errorf("expected FAQ section in body, got: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "post-author-card") || !strings.Contains(bodyStr, "Dafa Gareth") {
		t.Errorf("expected Author card in body, got: %s", bodyStr)
	}
}




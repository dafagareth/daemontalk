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
	if !strings.Contains(bodyStr, "fig-count") {
		t.Errorf("expected fig-count in body")
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
}

func TestDescriptionUnescape(t *testing.T) {
	src := []byte(`---
title: "Golang Concepts"
slug: golang-concepts
---

Sebagian besar tutorial Go berhenti di "cara mencetak Hello World", padahal masalah nyata & solusi baru muncul.
`)
	p, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	expected := `Sebagian besar tutorial Go berhenti di "cara mencetak Hello World", padahal masalah nyata & solusi baru muncul.`
	if p.Description != expected {
		t.Errorf("got %q, want %q", p.Description, expected)
	}
	if strings.Contains(p.Description, "&quot;") || strings.Contains(p.Description, "&amp;") {
		t.Errorf("description contains raw HTML entities: %q", p.Description)
	}
}

func TestDescriptionFrontmatter(t *testing.T) {
	src := []byte(`---
title: "Golang Concepts"
slug: golang-concepts
description: "Deskripsi khusus dari frontmatter & \"quotes\""
---

Paragraf pertama yang seharusnya di-override.
`)
	p, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	expected := `Deskripsi khusus dari frontmatter & "quotes"`
	if p.Description != expected {
		t.Errorf("got %q, want %q", p.Description, expected)
	}
}

func TestReferencesAndCalloutsParsing(t *testing.T) {
	content := `---
title: "Advanced Systems"
slug: advanced-systems
---

> [!WARNING]
> Pastikan mutex terkunci sebelum mengakses map!

` + "```callout TIP\n" + "Gunakan sync.Pool untuk mengurangi beban alokasi memori GC.\n```\n\n" + "```references\n" + "- title: The Linux Programming Interface\n  author: Michael Kerrisk\n  year: 2010\n  url: https://man7.org/tlpi/\n```\n"

	p, err := Parse([]byte(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	body := string(p.Body)

	// Verify GitHub-style Alert Warning
	if !strings.Contains(body, "post-callout") || !strings.Contains(body, "WARNING") || !strings.Contains(body, "Pastikan mutex terkunci") {
		t.Error("GitHub alert warning not parsed properly")
	}

	// Verify Code-fence Callout Tip
	if !strings.Contains(body, "TIP") || !strings.Contains(body, "sync.Pool") {
		t.Error("Callout TIP not parsed properly")
	}

	// Verify References & Citations block
	if !strings.Contains(body, "post-references") || !strings.Contains(body, "The Linux Programming Interface") || !strings.Contains(body, "Michael Kerrisk") || !strings.Contains(body, "man7.org/tlpi/") {
		t.Error("References block not parsed properly")
	}
}

func TestTabsLinkAndStatParsing(t *testing.T) {
	content := `---
title: "Advanced Systems Components"
slug: advanced-systems-components
---

` + "```tabs\n=== main.go\npackage main\nfunc main() {}\n\n=== Dockerfile\nFROM alpine\n```\n\n" + "```link\nurl: https://kernel.org\ntitle: The Linux Kernel Archives\ndescription: Primary site for the Linux kernel source code.\nsite: kernel.org\n```\n\n" + "```stat\n- value: \"14.8x\"\n  label: \"Throughput Boost\"\n  description: \"vs baseline sync.Mutex\"\n- value: \"0.8 µs\"\n  label: \"P99 Latency\"\n```\n"

	p, err := Parse([]byte(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	body := string(p.Body)

	// Verify Code Tabs
	if !strings.Contains(body, "data-code-tabs") || !strings.Contains(body, "main.go") || !strings.Contains(body, "Dockerfile") {
		t.Error("Tabs block not parsed properly")
	}

	// Verify Link Card
	if !strings.Contains(body, "https://kernel.org") || !strings.Contains(body, "The Linux Kernel Archives") || !strings.Contains(body, "kernel.org") {
		t.Error("Link card block not parsed properly")
	}

	// Verify Stat Grid
	if !strings.Contains(body, "post-stat-grid") || !strings.Contains(body, "14.8x") || !strings.Contains(body, "Throughput Boost") || !strings.Contains(body, "0.8 µs") {
		t.Error("Stat block not parsed properly")
	}
}

func TestLaTeXMathProtection(t *testing.T) {
	content := `---
title: "KV-Cache Math"
slug: kv-cache-math
---

$$\text{KV-Cache Per Request} = 2 \times n_{\text{layers}} \times n_{\text{heads}} \times d_{\text{head}} \times \text{seq\_len} \times \text{bytes}$$

Mari kita uji inline math: $n_{\text{layers}}$ dan $\text{seq\_len}$.
`

	p, err := Parse([]byte(content))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	body := string(p.Body)

	if strings.Contains(body, "xDTUSCOREx") || strings.Contains(body, "xDTESCAPEDUSCOREx") || strings.Contains(body, "xDTASTx") {
		t.Errorf("Leaked internal math placeholders in body: %s", body)
	}

	if !strings.Contains(body, `n_{\text{layers}}`) {
		t.Errorf("Subscript n_{\\text{layers}} was mangled: %s", body)
	}
	if !strings.Contains(body, `\text{seq\_len}`) {
		t.Errorf("Escaped seq\\_len was mangled: %s", body)
	}
}

func TestPostStatusParsing(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		wantStatus string
		wantDraft  bool
		wantType   string
	}{
		{
			name: "explicit status published",
			yaml: `---
title: "Published Article"
slug: pub-1
status: published
type: article
---
Body text.`,
			wantStatus: "published",
			wantDraft:  false,
			wantType:   "article",
		},
		{
			name: "explicit status draft",
			yaml: `---
title: "Draft Article"
slug: draft-1
status: draft
---
Body text.`,
			wantStatus: "draft",
			wantDraft:  true,
			wantType:   "article",
		},
		{
			name: "legacy draft false",
			yaml: `---
title: "Legacy Published"
slug: leg-pub
draft: false
---
Body text.`,
			wantStatus: "published",
			wantDraft:  false,
			wantType:   "article",
		},
		{
			name: "legacy draft true",
			yaml: `---
title: "Legacy Draft"
slug: leg-draft
draft: true
---
Body text.`,
			wantStatus: "draft",
			wantDraft:  true,
			wantType:   "article",
		},
		{
			name: "default without status or draft",
			yaml: `---
title: "Default Article"
slug: def-1
---
Body text.`,
			wantStatus: "published",
			wantDraft:  false,
			wantType:   "article",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Parse([]byte(tc.yaml))
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if p.Status != tc.wantStatus {
				t.Errorf("Status = %q, want %q", p.Status, tc.wantStatus)
			}
			if p.Draft != tc.wantDraft {
				t.Errorf("Draft = %v, want %v", p.Draft, tc.wantDraft)
			}
			if p.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", p.Type, tc.wantType)
			}
		})
	}
}

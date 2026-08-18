---
# Daemontalk Post Template
# Key frontmatter notes:
# - slug: Unique URL route (/blog/{slug}). Alphanumeric text or an 8-byte hex UID.
# - date: Publication date in YYYY-MM-DD format.
# - lang: Language code ('en' for English, 'id' for Indonesian).
# - draft: Set to 'true' for preview/draft mode; 'false' to publish.
# - description: 1-2 sentence summary for SEO and social media previews.

title: "Article Title Here"
slug: article-slug-here
aliases: []
date: 2026-08-17
author: "Author Name"
tags: ["topic-one", "topic-two"]
lang: en
draft: false
type: post
cover: "/static/images/posts/article-slug/cover.png"
coverCaption: "Cover image caption or attribution"
coverSource: "https://example.com"
readTime: 5
description: "A concise summary of the article for search engines and social previews."
series: "Optional Series Name"
series_part: 1
---

Write the opening paragraph here. This introduces the topic and serves as the article's executive summary.

## 1. Typography & Text Formatting

Demonstrate core text formatting and inline elements:

- **Bold text** for emphasis and key concepts
- *Italic text* for foreign terms or definitions
- ~~Strikethrough~~ for deprecated content
- `Inline code` for commands, variables, or functions
- [External links](https://example.com) for references
- Footnote citations[^1]

## Footnotes

[^1]: Additional context or citation details.

> A key quote, architectural principle, or notable insight.

## 2. Callout & Alert Blocks

> [!NOTE]
> Relevant background context or supplementary details.

> [!TIP]
> Practical advice or recommended best practices.

> [!IMPORTANT]
> Critical requirements or essential steps.

> [!WARNING]
> Warnings about edge cases or potential pitfalls.

> [!CAUTION]
> Warnings about high-risk or destructive actions.

## 3. Key Metrics & Statistics

```stat
- value: "10x"
  label: "Throughput"
  description: "Improvement under high concurrency"

- value: "99.9%"
  label: "Uptime"
  description: "Measured over 12 months in production"

- value: "-45%"
  label: "Memory Usage"
  description: "Reduction in heap allocations"
```

## 4. Multi-File Code Tabs

```tabs
=== main.go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Daemontalk!")
}

=== config.yaml
server:
  port: 8080
  host: "0.0.0.0"

=== Makefile
build:
	go build -o app .
```

## 5. Annotated Code Snippets

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
	conn := acquireConnection() // [!code hl]
	defer releaseConnection(conn)

	legacyHandler(w, r) // [!code --]
	modernHandler(w, r) // [!code ++]
}
```

## 6. Rich Link Preview

```link
url: https://example.com
title: Example Resource Title
description: A brief summary of the external resource and why it is relevant.
site: example.com
```

## 7. Interactive Checklist

- [ ] Task item one
- [ ] Task item two
- [ ] Task item three

## 8. Architecture Diagrams

```text
+-------------------+      +-------------------+
|    Client / UI    | ---> |    Backend API    |
+-------------------+      +-------------------+
                                     |
                                     v
+-------------------+      +-------------------+
|      Cache        | <--- |     Database      |
+-------------------+      +-------------------+
```

## 9. Comparison Matrix

| Feature | Solution A | Solution B | Solution C |
| :--- | :--- | :--- | :--- |
| Performance | High | Moderate | Low |
| Complexity | Low | Medium | High |
| Scalability | Horizontal | Vertical | Limited |

## 10. Media Components

### Carousel

```carousel
![Slide 1](/static/images/posts/example/slide-1.png "First slide description")
![Slide 2](/static/images/posts/example/slide-2.png "Second slide description")
```

### Gallery

```gallery
![Image 1](/static/images/posts/example/img-1.png "Gallery image one")
![Image 2](/static/images/posts/example/img-2.png "Gallery image two")
```

## 11. Frequently Asked Questions

```faq
Q: What is the primary use case for this approach?
A: It is ideal for systems requiring high throughput and predictable latency.

Q: How can this be tested in a local environment?
A: You can run the provided Makefile targets or execute the test suite directly.
```

## 12. Structured References

```references
- title: The Linux Programming Interface
  author: Michael Kerrisk
  year: 2010
  publisher: No Starch Press
  url: https://man7.org/tlpi/

- title: Designing Data-Intensive Applications
  author: Martin Kleppmann
  year: 2017
  publisher: O'Reilly Media
  url: https://dataintensive.net/
```

## 13. Author Card

```author
name: daemontalk team
role: Technical Writer
avatar: /static/logo/icon-dark.png
bio: Technical collective focusing on Linux systems, Go backend architecture, and distributed engineering.
github: https://github.com/dafagareth/daemontalk
email: team@daemontalk.dev
website: https://daemontalk.com
```

## 14. Mathematics (LaTeX)

Daemontalk supports academic-standard LaTeX formatting via KaTeX. 

You can write inline math using single dollar signs:
The time complexity is $O(n \log n)$ and the space complexity is $O(n)$.

For block math, use double dollar signs on their own lines or inline:
$$
f(x) = \int_{-\infty}^{\infty} \hat{f}(\xi)\,e^{2 \pi i \xi x} \, d\xi
$$

Or inline block: $$\text{KV-Cache} = 2 \times n_{\text{layers}} \times d_{\text{head}}$$
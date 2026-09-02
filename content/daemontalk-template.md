---
# Daemontalk Dispatch Template
# - slug: Unique URL route (/blog/{slug}). Use lowercase letters and hyphens.
# - author_github: GitHub username for automatic avatar resolution.
# - contributors: List of GitHub usernames who contributed edits or fixes.
# - type: Format badge ('dispatch' or 'article').
# - status: 'published' to go live; 'draft' for local preview.
# - readTime: Estimated reading duration in minutes.

title: "Dissecting Linux OOM Killer: How the Kernel Selects Memory Victims"
slug: "os-oom-killer-internals"
aliases: []
date: 2026-08-30
author: "daemontalk team"
author_github: "dafagareth"
contributors: []
tags: ["os", "linux", "systems", "kernel"]
lang: "en"
status: "published"
type: "dispatch"
readTime: 6
cover: "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Linux Kernel Memory Management"
coverSource: "https://kernel.org"
description: "An in-depth architectural breakdown of how the Linux Out-Of-Memory (OOM) Killer calculates oom_score and chooses processes to terminate."
series: "Linux Kernel Internals"
series_part: 1
---

<span class="drop-cap">T</span>HE opening paragraph should dive straight into the technical problem without conversational fluff. State the architectural bottleneck, performance trade-off, or kernel subsystem being dissected.

## 1. Editorial Drop Caps

Add an editorial drop cap to your opening letter or word:

- **Single Letter Drop Cap**: `<span class="drop-cap">W</span>hen the system runs out of memory...`
- **Multi-Letter Drop Cap**: `<span class="drop-cap">THE</span> kernel invokes the OOM killer...`

## 2. Text Formatting & Cross-References

Use standard Markdown with Daemontalk typographic extensions:

- **Bold**: `**kernel space**` for critical terminology.
- *Italics*: `*borrow checker*` for foreign or runtime concepts.
- `Inline code`: `sys_sendfile` for syscalls, structs, and variables.
- [Cross-Reference Link](/blog/protocols-http3 "cross-ref"): Add `title="cross-ref"` for dashed internal dispatches.
- Footnotes[^1]: Secondary context that avoids cluttering the main text.

[^1]: Footnotes provide additional nuance without disrupting the reading flow.

## 3. Live Wikipedia Previews & Abbreviations

Daemontalk automatically renders interactive hover cards and tooltips:

- **Live Wikipedia Preview**: Any standard link to Wikipedia will automatically display an interactive hover card with the live thumbnail and summary extract fetched from the Wikipedia API:
  `[Transmission Control Protocol](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)` or `[Linux Kernel](https://en.wikipedia.org/wiki/Linux_kernel)`.
- **Semantic Abbreviations**: `<abbr title="Internet Engineering Task Force">IETF</abbr>` or `<abbr title="Head-of-Line">HoL</abbr>`.

## 4. Images, Captions & Text Wrapping

Embed images as full-width blocks, with captions, or wrapped around text:

### Full-Width Image with Caption
![Kernel memory allocation graph](https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80 "Figure 1: Page cache allocation and swap space utilization under high load")

### Text Wrapping (Float Left & Float Right)
To wrap text around a diagram, use `<img class="float-right">` or `<img class="float-left">`:

<img src="https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=600&q=80" alt="CPU Cache Coherence" class="float-right">

Text placed immediately after a float image will naturally wrap around it on desktop screens, and automatically collapse to a centered block on mobile viewports. This is ideal for architecture diagrams, profile icons, or compact memory maps that accompany explanatory paragraphs.

<div class="clear-both"></div>

## 5. Media Carousels & Galleries

Present related diagrams side-by-side or in an interactive slider:

### Image Carousel (`carousel`)
```carousel
![Slide 1: TCP Handshake](https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80 "1. Initial SYN-ACK exchange")
![Slide 2: TLS 1.3 Handshake](https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1200&q=80 "2. Key exchange with Perfect Forward Secrecy")
```

### Multi-Column Gallery (`gallery`)
```gallery
![L1 Cache](https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=800&q=80 "L1 Data & Instruction Cache")
![L2 Cache](https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=800&q=80 "L2 Unified Non-Inclusive Cache")
```

## 6. Callouts & Alert Blocks

Daemontalk supports GitHub-style alert callouts:

> [!NOTE]
> Background context or supplementary architectural notes.

> [!TIP]
> Practical optimization advice or production configuration tips.

> [!IMPORTANT]
> Invariants, critical prerequisites, or mandatory verification steps.

> [!WARNING]
> Potential pitfalls, race conditions, or memory leak hazards.

> [!CAUTION]
> High-risk actions that can cause kernel panics or data destruction.

## 7. Key Metrics & Statistics (`stat`)

Highlight benchmark gains and hardware metrics visually:

```stat
- value: "10x"
  label: "Throughput"
  description: "Transfer rate increase via zero-copy sendfile"

- value: "99.99%"
  label: "Uptime"
  description: "Measured over 12 months in production"

- value: "-45%"
  label: "Context Switches"
  description: "Reduction in user-to-kernel boundary transitions"
```

## 8. Multi-File Code Tabs (`tabs`)

Group multi-file implementations or configurations into a tabbed widget:

```tabs
=== main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Daemontalk zero-copy service")
	})
	http.ListenAndServe(":8080", nil)
}

=== config.yaml
server:
  host: "127.0.0.1"
  port: 8080
  workers: 8

=== Makefile
build:
	go build -trimpath -ldflags="-s -w" -o server .
test:
	go test -count=1 -race ./...
```

## 9. Code Annotations & Diffs

Annotate code blocks with line highlights, additions, and removals:

```go
func allocateBuffer(size int) []byte {
	buf := make([]byte, size) // [!code hl]

	legacySlowInit(buf)       // [!code --]
	zeroCopyDirectInit(buf)   // [!code ++]

	return buf
}
```

## 10. ASCII Architecture & Flow Diagrams

Construct clean ASCII flow diagrams inside `text` code blocks:

```text
+------------------+       +------------------+       +------------------+
|  Userspace App   | --->  |   Linux Kernel   | --->  | Physical Device  |
| (Nginx / Kafka)  |       | (sendfile / io)  |       |  (NIC / NVMe)    |
+------------------+       +------------------+       +------------------+
                                     |
                                     v [Zero-Copy Bypass]
                           +------------------+
                           |    Page Cache    |
                           +------------------+
```

## 11. Structured Comparison Tables

| Feature | HTTP/1.1 | HTTP/2 | HTTP/3 |
| :--- | :--- | :--- | :--- |
| **Transport** | TCP | TCP | UDP (QUIC) |
| **Handshake** | 1-RTT + TLS | 1-RTT + TLS | 0-RTT / 1-RTT |
| **HoL Blocking** | Application Level | Transport Level | None (Independent Streams) |

## 12. Mathematics (LaTeX & KaTeX)

Render mathematical formulas natively:

- **Inline math**: Computational complexity is $O(n \log n)$ with $O(1)$ auxiliary space.
- **Block math**:

$$
\text{oom\_score} = \frac{\text{points} \times 1000}{\text{total\_pages}} + \text{oom\_score\_adj}
$$

## 13. Rich Link Preview Cards (`link`)

```link
url: https://github.com/torvalds/linux/blob/master/mm/oom_kill.c
title: Linux Kernel mm/oom_kill.c Source Code
description: Official implementation of the Out-Of-Memory killer process scoring algorithm.
site: github.com
```

## 14. Frequently Asked Questions (`faq`)

```faq
Q: Can oom_score_adj be set to negative values?
A: Yes. The allowed range is from -1000 (never terminate) to +1000 (always terminate first).

Q: How can memory pressure be tested safely?
A: Configure cgroups v2 memory.max or run stress-ng in an isolated test container.
```

## 15. Structured References (`references`)

Conclude every dispatch with authoritative citations (RFCs, kernel source, books, papers):

```references
- title: "Linux Kernel Documentation: Memory Management Concepts"
  author: "The Linux Kernel Organization"
  url: "https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html"

- title: "The Linux Programming Interface"
  author: "Michael Kerrisk"
  year: 2010
  publisher: "No Starch Press"
  url: "https://man7.org/tlpi/"

- title: "RFC 9000: QUIC - A UDP-Based Multiplexed and Secure Transport"
  author: "IETF QUIC Working Group"
  year: 2021
  url: "https://datatracker.ietf.org/doc/html/rfc9000"
```
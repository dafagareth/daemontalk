# Daemontalk Writing Guide

Daemontalk is an open engineering notebook and tech portfolio focusing on Linux systems, Go backends, and low-level software development. This document outlines Markdown formatting conventions and how to share technical notes or improvements.

---

## Writing Approach

The focus here is on practical, reproducible technical notes. If you would like to share laboratory experiments, production debugging experiences, or university computer science explorations, your contributions are very welcome.

Getting straight to the problem, architecture, or code examples in the opening paragraphs is recommended. Keeping introductions concise helps readers quickly grasp the technical essence of the topic.

---

## Article Formats

There are three primary formats shared on Daemontalk:

### 1. Technical Deep-Dives
Detailed explorations of systems concepts such as Linux kernel mechanics, Go concurrency models, or storage engines. Set `type: post` in the frontmatter.

### 2. Incident & Debugging Notes (RCA)
Step-by-step reconstructions of real-world bugs, diagnostic timelines, and root cause resolutions. Set `type: post` with `rca` or `incident` in tags.

### 3. Today I Learned (TIL)
Bite-sized notes focused on a single finding that can be read in a few minutes, such as a handy CLI workflow, compiler flag, or syscall behavior. Set `type: til` in the frontmatter.

---

## Frontmatter Specification

Articles are stored in `content/posts/` as Markdown files with YAML frontmatter at the top:

```yaml
---
title: "Zero-Copy I/O with io_uring in Go"
slug: "7f8a9b1c"
date: "2026-08-08"
tags: ["linux", "go", "performance", "storage"]
lang: "en"
draft: false
type: "post"
summary: "Exploring asynchronous I/O batching and ring-buffer submissions using io_uring in Go."
---
```

Use `lang: "en"` for English articles and `lang: "id"` for Indonesian articles. Always specify the language identifier for code blocks (such as go, bash, or c) for proper syntax highlighting.

You can download the full starter template file here: [Download template.md](/download/template.md).

---

## Submission Workflow

### Via GitHub Pull Request
1. Fork and clone the repository: `git clone https://github.com/dafagareth/daemontalk`.
2. Create a new branch: `git checkout -b post/your-topic-slug`.
3. Add your Markdown file to `content/posts/your-topic-slug.md`.
4. Test the build locally with `make build`.
5. Open a Pull Request on GitHub with a brief description of the topic.

### Via Email
If you prefer not using GitHub, you can send your standalone Markdown draft or git patch directly to **realdaemontalk@gmail.com** with the article title in the subject line.

---

## Licensing & Authorship

All written articles on Daemontalk are shared under the Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0) license, while code snippets are under the MIT License. Authors retain full attribution for their original work.

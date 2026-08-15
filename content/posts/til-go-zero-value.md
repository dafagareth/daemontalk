---
title: "TIL: Go's zero value makes struct initialization safe by default"
slug: b3b246ab
aliases: [til-go-zero-value]
date: 2026-06-06
tags: [go, til]
lang: en
draft: false
type: til
---

In Go, every variable is initialized to its **zero value** automatically: `0` for integers, `""` for strings, `nil` for pointers and maps, `false` for bools.

This means you never have to worry about uninitialized memory. A freshly allocated struct is always in a valid, usable state:

```go
var mu sync.Mutex   // ready to use, no New() needed
var buf bytes.Buffer // ready to use, no New() needed
var wg sync.WaitGroup
```

The standard library is designed around this: `sync.Mutex`, `bytes.Buffer`, `sync.WaitGroup` all work correctly straight out of the box without any constructor call.

**Why this matters:** In languages without zero values, forgetting to initialize a variable is a common source of bugs. Go eliminates that whole class of errors by design.

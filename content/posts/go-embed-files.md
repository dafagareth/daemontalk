---
title: "Embedding Files in Go Binaries with go:embed"
slug: 4fb361af
aliases: [go-embed-files]
date: 2025-05-20
tags: [go, tools, backend]
lang: en
draft: false
---

Before Go 1.16, shipping an application that needed static files meant bundling them separately, using a third-party tool to bake them into the binary, or reading from the filesystem at runtime. None of those options are clean.

`go:embed` solved this. It is a compiler directive that includes files or directories directly in the binary at build time. The result is a single self-contained executable.

## Basic Usage

```go
package main

import (
    _ "embed"
    "fmt"
)

//go:embed hello.txt
var content string

func main() {
    fmt.Println(content)
}
```

Place `hello.txt` in the same directory as your Go file. The compiler reads it at build time and stores it in `content`. No file is needed at runtime.

The `_ "embed"` blank import is required even if you call no function from the package. It activates the directive.

## Bytes Instead of Strings

```go
//go:embed logo.png
var logo []byte
```

Use `[]byte` for binary files, `string` for text.

## Embedding Directories with fs.FS

For multiple files or whole directories:

```go
import "embed"

//go:embed static
var staticFiles embed.FS
```

Access files through the `fs.FS` interface:

```go
data, err := staticFiles.ReadFile("static/css/main.css")
```

## Serving Over HTTP

The standard library's `http.FileServer` accepts an `fs.FS` directly:

```go
import (
    "embed"
    "io/fs"
    "net/http"
)

//go:embed static
var staticFiles embed.FS

func main() {
    sub, _ := fs.Sub(staticFiles, "static")
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
    http.ListenAndServe(":8080", nil)
}
```

`fs.Sub` strips the `static/` prefix so HTTP paths match file names directly.

## Glob Patterns

```go
//go:embed templates/*.html
var templates embed.FS

//go:embed static/css/*.css static/js/*.js
var assets embed.FS
```

## Hidden Files

By default, files starting with `.` are excluded. To include them:

```go
//go:embed all:static
var staticFiles embed.FS
```

## When to Use It

Embed files that do not change between deployments: CSS, fonts, HTML templates, SQL migration files. Do not embed user-uploaded content, per-environment configuration, or anything that changes without a rebuild.

For a deployment model where you want to copy one binary to a server and run it, `go:embed` removes the entire "where are my static files" problem.

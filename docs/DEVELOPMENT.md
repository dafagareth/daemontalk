# Development and Contribution Guide

This document outlines the local development workflow, build toolchain, testing procedures, and code standards for Daemontalk.

---

## 1. Prerequisites

Ensure the following tools are installed:
* Go 1.26 or higher
* Node.js 20 or higher
* `templ` CLI: `go install github.com/a-h/templ/cmd/templ@latest`

---

## 2. Build Pipeline

The Daemontalk build pipeline combines template generation, CSS minification, and Go binary compilation:

```bash
# Install dependencies
npm install

# Full build (generates templ files, compiles Tailwind CSS, and builds Go binary)
make build

# Start the application locally
./daemontalk
```

### Build Commands Reference

* `make templ`: Runs `templ generate` to compile `.templ` files to Go source files.
* `make css`: Compiles and minifies Tailwind CSS v4 input to `web/static/css/main.css`.
* `make build`: Executes full build sequence and produces the `daemontalk` binary.
* `make run`: Builds and executes the application on `http://localhost:8080` (HTTP) and `localhost:2222` (SSH).

---

## 3. Testing and Verification

Always verify all test suites before submitting pull requests or making modifications:

```bash
# Run all unit tests with race detector enabled
go test -race ./...

# Run clean test execution without caching
go test -count=1 ./...

# Run specific package tests
go test -v ./internal/auth/...
go test -v ./internal/forum/...
go test -v ./internal/handler/...
```

---

## 4. Code and Architecture Conventions

* **Server-Side Rendering**: HTML interfaces must use compiled `templ` components under `web/templates/`. Avoid introducing client-side JavaScript frameworks.
* **Database Access**: All SQL operations must use parameterized queries. Database logic belongs exclusively in `internal/<module>/store*.go`.
* **Zero-Trust Input Handling**: Never trust client-submitted user identifiers. User identity must be strictly extracted from `auth.GetUser(r.Context())`.
* **Prose Formatting in Documentation**: Editorial documents in `content/` and legal policies must use structured narrative paragraphs with bold lead-ins rather than bullet point lists.

---
title: "Value Canonicalization and String Interning with Go 1.23 unique Package"
slug: d8f2a9e5
aliases: [go-unique-package-interning]
date: 2026-07-18
tags: [go, performance, memory]
lang: en
draft: false
type: post
---

Large Go applications handling JSON payloads, HTTP headers, or database query results often allocate thousands of duplicate string instances across heap memory. Go 1.23 introduced the standard library `unique` package, providing a type-safe mechanism for value interning and canonicalization. By replacing duplicate allocations with globally interned handles, applications reduce memory footprint and enable fast pointer equality checks.

## Fun Facts

**Fact 1.** Go strings consist of a 16-byte header containing a pointer and a length integer. When duplicate strings exist in memory, each header points to an independent heap allocation even if the byte sequences are identical.

**Fact 2.** The `unique` package uses weak pointers inside the Go runtime garbage collector, allowing interned values to be automatically reclaimed when all external references to their `Handle[T]` expire.

**Fact 3.** Comparing two `unique.Handle[T]` values uses direct pointer comparison in a single CPU instruction, outperforming `bytes.Equal` or string byte-by-byte comparison on long text fields.

---

## Tips and Tricks

### 1. Canonicalize Strings Using unique.Make
Wrap string variables with `unique.Make` to retrieve a canonical `unique.Handle[string]`.

```go
package main

import (
	"fmt"
	"unique"
)

func main() {
	h1 := unique.Make("application/json")
	h2 := unique.Make("application/json")
	fmt.Println(h1 == h2) // Output: true
}
```

### 2. Intern Complex Comparable Structs
`unique.Make` accepts any comparable generic type, enabling canonicalization of metadata structs.

```go
type Header struct {
	Name, Value string
}

func InternHeader(name, val string) unique.Handle[Header] {
	return unique.Make(Header{Name: name, Value: val})
}
```

### 3. Compare unique.Make Against Manual sync.Map Interning
Traditional string interning using `sync.Map` causes memory leaks because stored entries remain rooted in global state. `unique.Make` uses GC weak references to prevent long-term memory accumulation.

```go
// Legacy sync.Map pattern (unbounded growth)
var globalCache sync.Map

func InternLegacy(s string) string {
	actual, _ := globalCache.LoadOrStore(s, s)
	return actual.(string)
}

// Go 1.23 pattern (GC aware)
func InternModern(s string) unique.Handle[string] {
	return unique.Make(s)
}
```

### 4. Benchmark Memory Footprint Reduction
Evaluate heap byte reduction by storing interned handles instead of raw repeated string instances.

```go
func BenchmarkMemory() {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	handles := make([]unique.Handle[string], 1000000)
	for i := range handles {
		handles[i] = unique.Make("tenant_identifier_string")
	}

	runtime.ReadMemStats(&m2)
	fmt.Printf("Allocated: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
}
```

### 5. Fast Symbol Table Lookups via Handle Equality
Use `unique.Handle` as map keys to accelerate lookup operations across large string pools.

```go
type SymbolTable struct {
	symbols map[unique.Handle[string]]int
}

func (st *SymbolTable) Register(token string) int {
	h := unique.Make(token)
	if id, exists := st.symbols[h]; exists {
		return id
	}
	id := len(st.symbols)
	st.symbols[h] = id
	return id
}
```

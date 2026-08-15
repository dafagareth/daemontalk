---
title: "Custom Dictionary Compression with Zstandard for Microservices"
slug: f1b4c6d9
aliases: [zstd-custom-dictionary-compression]
date: 2026-06-02
tags: [performance, tools, storage]
lang: en
draft: false
type: post
---

Standard compression algorithms like gzip or default Zstandard (zstd) underperform when applied to small API payloads under two kilobytes. Because small JSON or Protocol Buffer messages contain repetitive structural keys but insufficient body text to build LZ77 sliding windows on the fly, compression headers often increase payload size. Training and embedding a pre-shared zstd dictionary solves this overhead by providing pre-calculated entropy tables for microservices.

## Fun Facts

**Fact 1.** Zstandard dictionary training samples hundreds or thousands of representative JSON payloads to build a compact, 10-110 KB binary state file containing common key names, structural syntax, and frequent byte sequences.

**Fact 2.** Using a 64 KB zstd dictionary on 1 KB microservice responses typically achieves 3x to 5x higher compression ratios compared to un-dictioried zstd or gzip level 6.

**Fact 3.** The zstd dictionary format includes a 4-byte Magic Number (`0xEC30A437`) and a unique 32-bit Dictionary ID, enabling receivers to verify that they are using the matching dictionary version before decompression.

---

## Tips and Tricks

### 1. Train a Zstandard Dictionary from API Samples
Collect representative JSON responses into a folder and execute `zstd --train` to build a customized 64 KB dictionary file.

```bash
mkdir -p /tmp/samples
# Populate /tmp/samples with 1000+ representative JSON payloads
zstd --train /tmp/samples/* -o api_v1.dict --maxdict=64KB
```

### 2. Verify Dictionary Compression Ratios via CLI
Compare raw zstd compression against dictionary-assisted compression on a 500-byte API response file.

```bash
# Standard compression without dictionary
zstd -14 response.json -o response.zst
ls -lh response.json response.zst

# Dictionary-assisted compression
zstd -14 -D api_v1.dict response.json -o response_dict.zst
ls -lh response_dict.zst
```

### 3. Compress Payloads with Embedded Dictionary in Go
Embed the trained dictionary byte slice into your Go binary using `go:embed` and instantiate a reusable `klauspost/compress/zstd` writer with dictionary options.

```go
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/klauspost/compress/zstd"
)

//go:embed api_v1.dict
var dictBytes []byte

func CompressWithDict(src []byte) ([]byte, error) {
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderDict(dictBytes))
	if err != nil {
		return nil, err
	}
	return encoder.EncodeAll(src, make([]byte, 0, len(src))), nil
}
```

### 4. Decompress Payloads with Embedded Dictionary in Rust
Incorporate the dictionary in Rust microservices using the `zstd` crate and `include_bytes!` macro.

```rust
use std::io::Read;

const DICT_BYTES: &[u8] = include_bytes!("api_v1.dict");

pub fn decompress_payload(compressed: &[u8]) -> Result<Vec<u8>, std::io::Error> {
    let mut decoder = zstd::Decoder::with_dictionary(compressed, DICT_BYTES)?;
    let mut decompressed = Vec::new();
    decoder.read_to_end(&mut decompressed)?;
    Ok(decompressed)
}
```

### 5. Benchmark Compression Speed vs Payload Size
Measure microsecond latency overhead and compression ratios when passing 500B vs 50KB payloads through dictionary contexts.

```go
func BenchmarkDictionaryCompression(b *testing.B) {
	payload := []byte(`{"status":"success","user_id":100492,"roles":["admin","billing"],"events":[{"id":1,"type":"login"}]}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CompressWithDict(payload)
	}
}
```

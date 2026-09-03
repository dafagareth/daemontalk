## Topic Wishlist (*Call for Articles*)

If you are looking to write or contribute to Daemontalk, here is a curated list of high-priority topics sought by the community:

---

### 1. Operating Systems & Linux Kernel

- **`[WANTED]` Architectural Deep-Dive: `io_uring` vs `epoll` in High-Throughput Servers**: How kernel submission and completion ring buffers eliminate syscall overhead.
- **`[WANTED]` Memory Allocators Shootout: `jemalloc` vs `mimalloc` vs Go Runtime Allocator**: Fragmentation mitigation, thread-local caching, and arena allocation algorithms.
- **`[WANTED]` Kernel-Level Security Profiling with eBPF & Aya (Rust)**: Intercepting dangerous system calls with negligible performance overhead.

---

### 2. Concurrency & Language Runtimes

- **`[WANTED]` Lock-Free Data Structures: The Michael-Scott Queue**: Practical implementations of atomic Compare-And-Swap (CAS) algorithms without mutex lock contention.
- **`[WANTED]` Go Goroutine Work-Stealing Scheduler Internals**: How `sysmon`, the network poller, and M:N work-stealing balance multi-core CPU utilization.

---

### 3. Distributed Systems & Storage Engines

- **`[WANTED]` Google Spanner Architecture: TrueTime API and Atomic Clocks**: How GPS and atomic clocks solve global serializable transaction ordering.
- **`[WANTED]` Write-Ahead Logging (WAL) & Crash Recovery Internals**: Checkpointing mechanics, ARIES recovery protocols, and disk `fsync` cost analysis in SQLite and PostgreSQL.

---

### 4. Networking & Cryptography

- **`[WANTED]` Anatomy of BGP Hijacking & RPKI Cryptographic Route Validation**: How global internet routing is redirected and how cryptographic validation secures internet prefixes.
- **`[WANTED]` WireGuard Protocol Internals vs OpenVPN / IPsec**: Modern Curve25519 cryptography and in-kernel VPN performance.

---

## How to Claim a Topic

1. Download the starter template: `curl -s https://daemontalk.com/daemontalk-template.md -o your-topic.md`
2. Write your dispatch and submit a Pull Request to our [GitHub Repository](https://github.com/dafagareth/daemontalk).
3. Once merged, your GitHub profile is permanently attributed across the dispatch header and granted the **CONTRIBUTOR** badge.

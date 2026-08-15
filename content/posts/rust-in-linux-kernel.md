---
title: "Rust in the Linux Kernel"
slug: a3f7c921
aliases: rust-in-linux-kernel
date: 2026-05-10
tags: [rust, linux, security]
lang: en
draft: false
---

Linux 6.1, merged in December 2022, introduced Rust as a second supported language for kernel development. The decision followed years of debate, a formal proposal by Miguel Ojeda, and a 2021 vote among kernel developers that showed strong support. Rust is not replacing C; it is available as an opt-in language for new subsystems and drivers where memory safety guarantees are worth the added toolchain complexity.

## Fun Facts

**Fact 1.** The Rust for Linux project began as a Google-sponsored initiative. Google engineers submitted the first RFC patches in 2021 and have continued to fund Ojeda's full-time work on the integration.

**Fact 2.** The kernel enforces that Rust code must not use `unsafe` arbitrarily. All unsafe blocks must be documented with a `// SAFETY:` comment explaining the invariant that makes the operation sound.

**Fact 3.** As of kernel 6.8, the Apple M1/M2 GPU driver (`drm/asahi`) is one of the most substantial Rust subsystems in the tree, written almost entirely in Rust by Asahi Linux developers.

---

## Tips and Tricks

### 1. Check Whether Your Kernel Was Built with Rust Support

Before writing or loading any Rust kernel module, confirm that the running kernel includes Rust support.

```bash
# Check for Rust support in the running kernel config
zcat /proc/config.gz | grep CONFIG_RUST

# Or from the build directory
grep CONFIG_RUST /boot/config-$(uname -r)
```

You should see `CONFIG_RUST=y`. If it is not set, you need to build your own kernel or use a distribution that enables it (Fedora 38+ includes it).

### 2. Enable Rust in Kconfig When Building the Kernel

When configuring a custom kernel build, Rust support lives under a specific menu. The toolchain must be detected automatically.

```bash
# Inside your kernel source directory
make menuconfig
# Navigate to: General setup > Rust support

# Or set it directly
scripts/config --enable CONFIG_RUST

# Verify that the Rust toolchain is found
make LLVM=1 rustavailable
```

The `rustavailable` target checks that `rustc`, `bindgen`, and `rustfmt` are present and at compatible versions. The kernel requires a specific minimum `rustc` version pinned in `Documentation/process/changes.rst`.

### 3. Write a Minimal Rust Kernel Module

A Rust out-of-tree module follows the same Kbuild conventions as a C module, but uses a `rust/` subdirectory and the `kernel` crate.

```rust
// rust/hello_kernel.rs
#![no_std]
#![feature(allocator_api)]

use kernel::prelude::*;

module! {
    type: HelloKernel,
    name: "hello_kernel",
    author: "Example Author",
    description: "A minimal Rust kernel module",
    license: "GPL",
}

struct HelloKernel;

impl kernel::Module for HelloKernel {
    fn init(_name: &'static CStr, _module: &'static ThisModule) -> Result<Self> {
        pr_info!("Hello from Rust kernel module\n");
        Ok(HelloKernel)
    }
}

impl Drop for HelloKernel {
    fn drop(&mut self) {
        pr_info!("Goodbye from Rust kernel module\n");
    }
}
```

The corresponding `Kbuild` file:

```makefile
# Kbuild
obj-$(CONFIG_HELLO_KERNEL) += hello_kernel.o
```

Build it with `make -C /lib/modules/$(uname -r)/build M=$(pwd) modules`.

### 4. Inspect Rust Modules Already in the Tree

The kernel source tree organizes Rust code under `rust/` (core abstractions and bindings) and `drivers/` or `fs/` for actual subsystem code.

```bash
# List all .rs files in the upstream kernel source
find drivers/ -name "*.rs" | head -20

# Example: look at the null block driver written in Rust
ls drivers/block/rnull.rs

# Count lines of Rust in the kernel tree
find . -name "*.rs" -not -path "*/target/*" | xargs wc -l | tail -1
```

Notable Rust subsystems as of 6.8: `drivers/gpu/drm/asahi/` (Apple GPU), `drivers/block/rnull.rs` (null block device sample), and the `net/phy/` Rust PHY abstractions.

### 5. Understand the Memory Safety Argument

The Linux kernel Security Team has tracked that roughly 65-70% of kernel CVEs historically involve memory safety bugs: use-after-free, out-of-bounds writes, and uninitialized reads. Rust eliminates entire classes of these at compile time through its ownership and borrow checker model.

```bash
# Search NVD for recent kernel memory safety CVEs (requires nvd-cli or curl)
curl -s "https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=linux+kernel+use-after-free&resultsPerPage=5" \
  | python3 -m json.tool | grep -E '"id"|"descriptions"' | head -20
```

The practical implication: a driver written correctly in Rust cannot produce a use-after-free or a null pointer dereference from safe code. The compiler rejects such programs. This is not a theoretical claim but a consequence of the type system, verified by the Rust compiler at build time rather than by runtime sanitizers.

---
title: "Landlock LSM Unprivileged Application Sandboxing in Linux"
slug: a9b8c7d6
aliases: [linux-landlock-sandboxing]
date: 2026-06-08
tags: [linux, security, tools]
lang: en
draft: false
type: post
---

Landlock is a Linux Security Module (LSM) introduced in kernel 5.13 that enables unprivileged applications to restrict their own access to filesystem paths and network ports. Traditional access control security modules like SELinux or AppArmor require administrator privileges to load policy rules. This post covers Landlock API mechanics, minimal C implementation code, comparison with seccomp-bpf, and usage in OpenSSH.

## Fun Facts

**Fact 1.** Landlock was created by Mickaël Salaün and integrated into Linux kernel 5.13 to enable unprivileged process sandboxing without CAP_SYS_ADMIN capabilities.

**Fact 2.** Landlock rulesets apply transitively to all child processes spawned after calling the restrict_self system call, enforcing persistent isolation across process subtrees.

**Fact 3.** OpenSSH adopted Landlock sandboxing in version 9.0 to isolate unprivileged pre-authentication child processes handling untrusted client network sessions.

---

## Tips and Tricks

### 1. Enforce Path Restrictions with C System Calls

Implement filesystem access restrictions by creating a Landlock ruleset, populating path descriptors, and restricting the calling thread.

```c
#include <fcntl.h>
#include <linux/landlock.h>
#include <sys/prctl.h>
#include <sys/syscall.h>
#include <unistd.h>

static inline int landlock_create_ruleset(const struct landlock_ruleset_attr *attr, size_t size, uint32_t flags) {
    return syscall(__NR_landlock_create_ruleset, attr, size, flags);
}

static inline int landlock_add_rule(int ruleset_fd, enum landlock_rule_type type, const void *attr, uint32_t flags) {
    return syscall(__NR_landlock_add_rule, ruleset_fd, type, attr, flags);
}

static inline int landlock_restrict_self(int ruleset_fd, uint32_t flags) {
    return syscall(__NR_landlock_restrict_self, ruleset_fd, flags);
}
```

### 2. Configure Allowed Read and Write Path Rules

Initialize the ruleset attribute structure to define permissible read or write operations before applying the policy.

```c
int restrict_filesystem_access(const char *allowed_path) {
    struct landlock_ruleset_attr attr = {
        .handled_access_fs = LANDLOCK_ACCESS_FS_READ_FILE | LANDLOCK_ACCESS_FS_READ_DIR
    };
    int ruleset_fd = landlock_create_ruleset(&attr, sizeof(attr), 0);
    if (ruleset_fd < 0) return -1;

    int path_fd = open(allowed_path, O_PATH | O_CLOEXEC);
    struct landlock_path_beneath_attr path_attr = {
        .allowed_access = LANDLOCK_ACCESS_FS_READ_FILE | LANDLOCK_ACCESS_FS_READ_DIR,
        .parent_fd = path_fd
    };
    landlock_add_rule(ruleset_fd, LANDLOCK_RULE_PATH_BENEATH, &path_attr, 0);
    close(path_fd);

    prctl(PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0);
    landlock_restrict_self(ruleset_fd, 0);
    close(ruleset_fd);
    return 0;
}
```

### 3. Compare Landlock and Seccomp-BPF Mechanisms

Landlock operates alongside seccomp-bpf to provide defense-in-depth security layers.

| Security Layer | Primary Scope | Granularity | Capability Requirement |
|---|---|---|---|
| Landlock | Filesystem paths & network ports | Object level (inodes, ports) | Unprivileged |
| Seccomp-BPF | System call filtering | Interface level (syscall numbers) | PR_SET_NO_NEW_PRIVS |

### 4. Inspect Landlock Capabilities in OpenSSH

Verify Landlock LSM support in installed OpenSSH daemon builds.

```bash
# Check kernel Landlock support status
dmesg | grep -i landlock

# Query Landlock ABI version support using liblandlock test binaries
ssh -V
```

### 5. Combine Landlock Restrictions with Seccomp Filters

Use seccomp to restrict dangerous system calls like `execve` while using Landlock to isolate directory read access.

```bash
# Verify process status sandboxing flags
cat /proc/self/status | grep -i seccomp
```

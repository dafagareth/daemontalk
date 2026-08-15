---
title: "Jujutsu Version Control Architecture and Git Integration"
slug: b4c3d2e1
aliases: [jujutsu-vcs-git-compatible]
date: 2026-07-12
tags: [tools, rust, devops]
lang: en
draft: false
type: post
---

Jujutsu (jj) is a Git-compatible version control system written in Rust by Martin von Zweigbergk. It provides a distinct user interface and concurrency model while maintaining full compatibility with standard Git repository storage formats. This post covers Jujutsu's core architecture, including anonymous commit tracking, operation logs, concurrent conflict storage, and remote Git workflows.

## Fun Facts

**Fact 1.** Martin von Zweigbergk created Jujutsu at Google to combine the workflow safety of Mercurial with the storage model and ecosystem of Git.

**Fact 2.** Jujutsu eliminates the traditional staging area index. Every modification in the working copy automatically updates an ongoing commit snapshot in real time.

**Fact 3.** Conflicts in Jujutsu are recorded as first-class objects within commits, allowing rebase operations to finish cleanly without halting for immediate manual conflict resolution.

---

## Tips and Tricks

### 1. Initialize Jujutsu on an Existing Git Repository

You can initialize Jujutsu inside any standard Git repository workspace using the colocated repository option.

```bash
# Navigate to your existing Git project
cd /path/to/project

# Co-locate Jujutsu tracking within the .git directory
jj git init --colocated

# Verify active revisions and working copy state
jj status
```

### 2. Manage Changes with Anonymous Commits

In Jujutsu, commits do not require immediate branch names or tag labels. Revisions are tracked using unique stable Change IDs.

```bash
# View recent revision history graph
jj log

# Create a new revision on top of current working copy
jj new -m "Refactor parser module"

# Describe or update commit messages for current revision
jj describe -m "Implement JSON parser tokenizer"
```

### 3. Revert Operations using the Undo Operation Log

Jujutsu records every command and state mutation in an internal operation log, making all destructive operations reversible.

```bash
# Display recent operations log history
jj op log

# Revert the most recent version control operation
jj undo

# Restore repository state to a specific operation ID
jj op restore 4a7c91e8b23f
```

### 4. Work with First-Class Conflicts

When a conflict occurs during a rebase or merge, Jujutsu creates a conflict state inside the commit and lets you continue working.

```bash
# Rebase current revision tree onto main branch
jj rebase -s @ -d main

# Inspect files containing unresolved conflict markers
jj diff

# Resolve conflict in editor then complete revision update
jj squash
```

### 5. Fetch and Push to Git Remote Repositories

Jujutsu synchronizes directly with standard Git remotes without requiring external conversion tools.

```bash
# Fetch changes from upstream Git remotes
jj git fetch

# Track a remote branch in Jujutsu
jj branch track main@origin

# Push local change IDs to remote Git references
jj git push --change @
```

---
title: "Git Worktrees: Work on Multiple Branches Without Stashing"
slug: 6b944d2d
aliases: [git-worktrees]
date: 2025-04-12
tags: [git, workflow, tools]
lang: en
draft: false
---

You are mid-feature when a production bug comes in. The usual move: stash your changes, check out main, fix the bug, pop the stash, hope nothing conflicts. It works, but it is unnecessary friction.

Git worktrees remove that friction. A worktree is a separate working directory linked to the same repository. Each worktree can be on a different branch, and they all share one `.git` directory.

## Creating a Worktree

```bash
git worktree add ../my-repo-hotfix hotfix/login-crash
```

This creates `../my-repo-hotfix` with `hotfix/login-crash` checked out. Your original directory stays untouched, mid-feature and all. Open a second terminal, `cd` into the new directory, and fix the bug.

To create a worktree with a new branch:

```bash
git worktree add -b hotfix/login-crash ../my-repo-hotfix main
```

## Practical Uses

**Review a PR locally without abandoning your branch:**

```bash
git fetch origin pull/42/head:pr-42
git worktree add ../review-pr42 pr-42
```

Open `../review-pr42` in your editor. Your main workspace is untouched.

**Run tests on a stable release while continuing development:**

```bash
git worktree add ../stable-test v2.1.0
cd ../stable-test && go test ./...
```

## Listing and Removing Worktrees

```bash
git worktree list
```

```
/home/dafa/my-repo        a1b2c3d [main]
/home/dafa/my-repo-hotfix e4f5g6h [hotfix/login-crash]
```

When you are done:

```bash
git worktree remove ../my-repo-hotfix
```

If there are uncommitted changes, Git refuses. Pass `--force` only when you are certain you want to discard them.

## One Rule to Know

You cannot check out the same branch in two worktrees simultaneously. Git prevents it. If you need the same branch in two places, create a branch from it first:

```bash
git worktree add ../inspect-main review/main-snapshot main
```

Everything else behaves identically to your main working directory: staging, committing, pushing, pulling all work normally.

Worktrees are worth the overhead when the context switch is meaningful. For a two-line fix, stashing is faster. For anything that requires compiling a project or running a dev server in a separate context, a worktree pays off immediately.

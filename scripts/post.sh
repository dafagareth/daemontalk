#!/usr/bin/env bash
# DaemonTalk Technical Post & Content Management CLI Entrypoint
# Usage: ./scripts/post.sh [command] [args...]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="${SCRIPT_DIR}/lib"

# Load modular components
source "${LIB_DIR}/common.sh"
source "${LIB_DIR}/post_create.sh"
source "${LIB_DIR}/post_lifecycle.sh"
source "${LIB_DIR}/post_inspect.sh"

ACTION="${1:-help}"
shift || true

case "$ACTION" in
    new|create)        cmd_new "$@" ;;
    new-uid|uid|hex)   cmd_new --uid "$@" ;;
    list|ls)           cmd_list ;;
    publish|pub)       cmd_publish "$@" ;;
    draft|unpublish)   cmd_draft "$@" ;;
    archive)           cmd_archive "$@" ;;
    restore)           cmd_restore "$@" ;;
    delete|rm)         cmd_delete "$@" ;;
    stats)             cmd_stats ;;
    validate|check)    cmd_validate ;;
    help|--help|-h)
        echo "DaemonTalk Content Management CLI"
        echo "Usage: ./scripts/post.sh <command> [arguments]"
        echo ""
        echo "Available commands:"
        echo "  new [title]          Create a new markdown post with readable slug"
        echo "  new --uid [title]    Create a new markdown post with 8-char hex UID"
        echo "  uid [title]          Shortcut: create new post directly with hex UID"
        echo "  list                 List all articles with status and dates"
        echo "  publish <slug>       Mark post as published and update publication date"
        echo "  draft <slug>         Revert post to draft status"
        echo "  archive <slug>       Archive post (removes from public feed)"
        echo "  restore <slug>       Restore archived post"
        echo "  delete <slug>        Permanently delete post and asset directory"
        echo "  stats                Print content publication metrics"
        echo "  validate             Validate frontmatter integrity across all posts"
        ;;
    *)
        echo "[error] Unknown command: '$ACTION'. Run './scripts/post.sh help' for usage."
        exit 1
        ;;
esac

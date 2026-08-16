#!/usr/bin/env bash
# DaemonTalk Content Management - Inspection & Validation Module

set -euo pipefail

cmd_list() {
    echo "Slug                              Status     Date         Title"
    echo "--------------------------------------------------------------------------------"
    for f in "${POSTS_DIR}"/*.md "${POSTS_DIR}"/*.md.archive; do
        [ -f "$f" ] || continue
        local slug
        local title
        local date_val
        local is_draft="false"
        local is_arch="false"

        slug=$(grep "^slug:" "$f" | head -1 | sed -E 's/slug:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "-")
        title=$(grep "^title:" "$f" | head -1 | sed -E 's/title:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "-")
        date_val=$(grep "^date:" "$f" | head -1 | sed -E 's/date:[[:space:]]*([^ ]+)/\1/' || echo "-")
        
        if grep -q "^draft:[[:space:]]*true" "$f" 2>/dev/null; then
            is_draft="true"
        fi
        if [[ "$f" == *.md.archive ]]; then
            is_arch="true"
        fi

        local status="public"
        if [ "$is_arch" = "true" ]; then
            status="archived"
        elif [ "$is_draft" = "true" ]; then
            status="draft"
        fi

        printf "%-32s  %-9s  %-10s  %s\n" "${slug:0:32}" "$status" "$date_val" "${title:0:40}"
    done
}

cmd_stats() {
    local total=0
    local published=0
    local drafts=0
    local archived=0

    for f in "${POSTS_DIR}"/*.md "${POSTS_DIR}"/*.md.archive; do
        [ -f "$f" ] || continue
        total=$((total + 1))
        if [[ "$f" == *.md.archive ]]; then
            archived=$((archived + 1))
        elif grep -q "^draft:[[:space:]]*true" "$f" 2>/dev/null; then
            drafts=$((drafts + 1))
        else
            published=$((published + 1))
        fi
    done

    echo "Content Metrics Summary:"
    echo "  Total Dispatches : $total"
    echo "  Published        : $published"
    echo "  Drafts           : $drafts"
    echo "  Archived         : $archived"
}

cmd_validate() {
    local errors=0
    echo "[info] Running content frontmatter validation..."

    for f in "${POSTS_DIR}"/*.md; do
        [ -f "$f" ] || continue
        local fname
        fname=$(basename "$f")

        if ! grep -q "^title:" "$f"; then
            echo "[warn] $fname: Missing 'title' in frontmatter"
            errors=$((errors + 1))
        fi
        if ! grep -q "^slug:" "$f"; then
            echo "[warn] $fname: Missing 'slug' in frontmatter"
            errors=$((errors + 1))
        fi
        if ! grep -q "^date:" "$f"; then
            echo "[warn] $fname: Missing 'date' in frontmatter"
            errors=$((errors + 1))
        fi
    done

    if [ "$errors" -eq 0 ]; then
        echo "[ok] All posts passed frontmatter validation."
    else
        echo "[error] Validation completed with $errors warning(s)."
    fi
}

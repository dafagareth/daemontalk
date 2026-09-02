#!/usr/bin/env bash
# DaemonTalk Content Management - Inspection & Validation Module

set -euo pipefail

cmd_list() {
    echo "Slug                              Status     Date         Type     Title"
    echo "-------------------------------------------------------------------------------------------------"
    for f in "${POSTS_DIR}"/*.md "${POSTS_DIR}"/*.md.archive; do
        [ -f "$f" ] || continue
        local slug
        local title
        local date_val
        local type_val
        local status_val=""
        local is_arch="false"

        slug=$(grep "^slug:" "$f" | head -1 | sed -E 's/slug:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "-")
        title=$(grep "^title:" "$f" | head -1 | sed -E 's/title:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "-")
        date_val=$(grep "^date:" "$f" | head -1 | sed -E 's/date:[[:space:]]*([^ ]+)/\1/' || echo "-")
        type_val=$(grep "^type:" "$f" | head -1 | sed -E 's/type:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "article")
        status_val=$(grep "^status:" "$f" | head -1 | sed -E 's/status:[[:space:]]*["'\'']?([^"'\'']+)["'\'']?/\1/' || echo "")
        
        if [ -z "$status_val" ]; then
            if grep -q "^draft:[[:space:]]*true" "$f" 2>/dev/null; then
                status_val="draft"
            else
                status_val="published"
            fi
        fi

        if [[ "$f" == *.md.archive ]]; then
            status_val="archived"
        fi

        printf "%-32s  %-9s  %-10s  %-7s  %s\n" "${slug:0:32}" "$status_val" "$date_val" "${type_val:0:7}" "${title:0:36}"
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
        elif grep -q -E "^status:[[:space:]]*['\"]?draft['\"]?" "$f" 2>/dev/null || grep -q "^draft:[[:space:]]*true" "$f" 2>/dev/null; then
            drafts=$((drafts + 1))
        else
            published=$((published + 1))
        fi
    done

    echo "Content Metrics Summary:"
    echo "  Total Articles   : $total"
    echo "  Published        : $published"
    echo "  Drafts           : $drafts"
    echo "  Archived         : $archived"
}

cmd_validate() {
    local errors=0
    local warnings=0
    echo "[info] Running content frontmatter validation across all articles..."

    for f in "${POSTS_DIR}"/*.md; do
        [ -f "$f" ] || continue
        local fname
        fname=$(basename "$f")

        if ! grep -q "^title:" "$f"; then
            echo "[error] $fname: Missing 'title' in frontmatter"
            errors=$((errors + 1))
        fi
        if ! grep -q "^slug:" "$f"; then
            echo "[error] $fname: Missing 'slug' in frontmatter"
            errors=$((errors + 1))
        fi
        if ! grep -q "^date:" "$f"; then
            echo "[error] $fname: Missing 'date' in frontmatter"
            errors=$((errors + 1))
        fi
        if ! grep -q "^type:" "$f"; then
            echo "[warn]  $fname: Missing 'type' in frontmatter (defaults to 'article')"
            warnings=$((warnings + 1))
        fi
        if ! grep -q -E "^(status|draft):" "$f"; then
            echo "[warn]  $fname: Missing 'status' (recommended: 'status: published' or 'status: draft')"
            warnings=$((warnings + 1))
        fi
        if ! grep -q "^lang:" "$f"; then
            echo "[warn]  $fname: Missing 'lang' in frontmatter"
            warnings=$((warnings + 1))
        fi
        if ! grep -q "^author:" "$f"; then
            echo "[warn]  $fname: Missing 'author' in frontmatter"
            warnings=$((warnings + 1))
        fi
    done

    if [ "$errors" -eq 0 ] && [ "$warnings" -eq 0 ]; then
        echo "[ok] All articles passed full frontmatter validation."
    elif [ "$errors" -eq 0 ]; then
        echo "[ok] Frontmatter valid with $warnings minor warning(s)."
    else
        echo "[error] Validation completed with $errors error(s) and $warnings warning(s)."
        return 1
    fi
}

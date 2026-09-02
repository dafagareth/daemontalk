#!/usr/bin/env bash
# DaemonTalk Content Management - Post Creation Module

set -euo pipefail

generate_hex_uid() {
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -hex 4
    elif [ -c /dev/urandom ] && command -v od >/dev/null 2>&1; then
        od -vAn -N4 -tx1 /dev/urandom | tr -d ' \n'
    elif [ -c /dev/urandom ] && command -v hexdump >/dev/null 2>&1; then
        hexdump -n 4 -e '4/1 "%02x"' /dev/urandom
    else
        date +%s%N | sha256sum | head -c 8
    fi
}

cmd_new() {
    local use_uid="false"
    local title=""
    local custom_slug=""

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --uid|-u|--hex|-h)
                use_uid="true"
                shift
                ;;
            *)
                if [ -z "$title" ]; then
                    title="$1"
                elif [ -z "$custom_slug" ]; then
                    custom_slug="$1"
                fi
                shift
                ;;
        esac
    done

    if [ -z "$title" ]; then
        read -p "Enter post title: " -r title
    fi
    if [ -z "$title" ]; then
        echo "[error] Title cannot be empty."
        exit 1
    fi

    local slug
    local readable_slug
    readable_slug=$(generate_slug "$title")

    if [ "$use_uid" = "true" ]; then
        # Generate 8-character hexadecimal UID (e.g. 7f4a9b2c)
        slug=$(generate_hex_uid)
    elif [ -n "$custom_slug" ]; then
        slug="$custom_slug"
    else
        slug="$readable_slug"
    fi

    local filename="${POSTS_DIR}/${slug}.md"
    local img_dir="${IMG_BASE_DIR}/${slug}"

    if [ -f "$filename" ]; then
        slug="${slug}-$(generate_hex_uid | cut -c 1-4)"
        filename="${POSTS_DIR}/${slug}.md"
        img_dir="${IMG_BASE_DIR}/${slug}"
    fi

    mkdir -p "${POSTS_DIR}" "${img_dir}"
    local today
    today=$(date +%Y-%m-%d)

    local aliases_str="[]"
    if [ "$use_uid" = "true" ] && [ -n "$readable_slug" ] && [ "$slug" != "$readable_slug" ]; then
        aliases_str="[\"${readable_slug}\"]"
    fi

    cat <<EOF > "$filename"
---
title: "${title}"
slug: "${slug}"
aliases: ${aliases_str}
date: ${today}
author: "daemontalk team"
tags: ["tag1", "tag2"]
lang: "en"
type: article
status: draft
description: "Write a brief summary of this technical article."
cover: "/static/images/posts/${slug}/cover.webp"
coverCaption: "Cover illustration description"
coverSource: "https://unsplash.com"
readTime: 5
---

Write the technical article body here using GitHub Flavored Markdown.
EOF

    echo "[ok] Post draft created: $filename"
    echo "[ok] Slug: $slug"
    if [ "$aliases_str" != "[]" ]; then
        echo "[ok] Human Alias (Auto-redirected): $aliases_str"
    fi
    echo "[ok] Asset folder created: $img_dir"

    if [ -n "${EDITOR:-}" ]; then
        "${EDITOR}" "$filename"
    fi
}

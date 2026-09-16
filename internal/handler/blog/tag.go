package blog

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	blogtmpl "daemontalk/web/templates/blog"
	"daemontalk/web/templates/layouts"
	portaltmpl "daemontalk/web/templates/portal"
	"daemontalk/web/templates/shared"

	"github.com/go-chi/chi/v5"
)

// PillarTags maps the 5 universal pillars to their respective child tags.
var PillarTags = map[string][]string{
	"wire": {
		"wire", "linux", "kernel", "networking", "security", "storage",
		"devops", "docker", "ebpf", "io-uring", "systemd", "sysadmin",
		"wireguard", "cgroups", "caddy", "nginx", "network", "ssh", "ufw",
		"openssl", "btrfs", "zfs", "landlock", "tailscale", "podman",
	},
	"craft": {
		"craft", "go", "rust", "backend", "database", "sqlite",
		"architecture", "concurrency", "performance", "algorithms", "python",
		"software", "api", "refactoring", "testing", "duckdb", "clickhouse",
	},
	"tools": {
		"tools", "cli", "terminal", "neovim", "tmux", "git", "bash",
		"fzf", "jq", "debugging", "workflow", "productivity", "shell",
		"editor", "ripgrep", "helix", "htop", "makefile", "zsh", "curl",
		"ffmpeg", "strace", "dotfiles", "tar", "zstd", "jujutsu",
	},
	"radar": {
		"radar", "ai", "llm", "agents", "open-weight", "inference",
		"wasm", "quantization", "microvm", "future", "science", "qemu",
		"firecracker", "dpo",
	},
	"essays": {
		"essays", "essay", "opinion", "career", "freelancing", "culture",
		"critique", "industry", "manifesto", "dispatch", "freelance",
		"education", "salary",
	},
}

func matchPostTag(postTags []string, tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if children, ok := PillarTags[tag]; ok {
		for _, pt := range postTags {
			ptLower := strings.ToLower(pt)
			for _, ct := range children {
				if ptLower == ct {
					return true
				}
			}
		}
		return false
	}
	for _, t := range postTags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

func (h *Handler) TagIndex(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	tag := chi.URLParam(r, "tag")
	isAdmin := h.isAdmin(r)

	var filtered []post.Post
	for _, p := range h.getVisiblePosts(isAdmin) {
		if matchPostTag(p.Tags, tag) {
			filtered = append(filtered, p)
		}
	}

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			slog.Error("tag view counts query failed", "tag", tag, "error", err)
		} else {
			viewCounts = vc
		}
	}

	h.render(w, r, layouts.Layout(ui, lang, "#"+tag, r.URL.Path, shared.PageMeta{
		Description: fmt.Sprintf("Posts tagged #%s on daemontalk.com", tag),
	}, blogtmpl.TagPage(ui, tag, filtered, lang, viewCounts)))
}

func (h *Handler) TagPostsPartial(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	tag := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag")))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 8 {
		offset = 8
	}

	isAdmin := h.isAdmin(r)
	var filtered []post.Post
	for _, p := range h.getVisiblePosts(isAdmin) {
		if matchPostTag(p.Tags, tag) {
			filtered = append(filtered, p)
		}
	}

	total := len(filtered)
	if offset >= total {
		h.render(w, r, portaltmpl.TagRiverItems(ui, nil, lang, offset, 0, 0, tag))
		return
	}

	batchSize := 12
	end := offset + batchSize
	if end > total {
		end = total
	}

	pagePosts := filtered[offset:end]
	nextOffset := end
	remaining := total - end
	h.render(w, r, portaltmpl.TagRiverItems(ui, pagePosts, lang, offset, nextOffset, remaining, tag))
}

func (h *Handler) RedirectTag(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	if tag != "" {
		http.Redirect(w, r, "/blog/tag/"+url.PathEscape(strings.ToLower(tag)), http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

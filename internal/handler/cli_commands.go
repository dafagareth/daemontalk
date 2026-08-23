package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CLIDaily(w http.ResponseWriter, r *http.Request) {
	if !IsCLIRequest(r) {
		lang := langFromRequest(r)
		ui := i18n.Get(lang)
		posts := h.AllPosts()

		var articles []post.Post
		for _, p := range posts {
			if p.Draft {
				continue
			}
			if len(articles) < 8 {
				articles = append(articles, p)
			}
		}

		meta := templates.PageMeta{
			Description: "Daily executive technical briefing and terminal CLI endpoints for systems engineers.",
		}
		h.Render(w, r, templates.Layout(ui, lang, "daily", r.URL.Path, meta,
			templates.DailyPage(ui, lang, articles, r.Host)))
		return
	}

	color := r.URL.Query().Get("nocolor") != "1"
	posts := h.AllPosts()

	var b strings.Builder
	banner := `  ____                                _        _ _    
 |  _ \  __ _  ___ _ __ ___   ___  _ __ | |_ __ _| | | __
 | | | |/ _` + "`" + ` |/ _ \ '_ ` + "`" + ` _ \ / _ \| '_ \| __/ _` + "`" + ` | | |/ /
 | |_| | (_| |  __/ | | | | | (_) | | | | || (_| | |   < 
 |____/ \__,_|\___|_| |_| |_|\___/|_| |_|\__\__,_|_|_|\_\`

	b.WriteString(fmt.Sprintf("%s%s%s\n\n", ansiGreen+ansiBold, banner, ansiReset))
	b.WriteString(fmt.Sprintf(" %s[ DAEMONTALK DAILY TECH BRIEFING ]%s\n", ansiCyan+ansiBold, ansiReset))
	b.WriteString(fmt.Sprintf(" %sDate:%s %s | %sHost:%s %s\n", ansiDim, ansiReset, time.Now().Format("Monday, 02 January 2006"), ansiDim, ansiReset, r.Host))
	b.WriteString(fmt.Sprintf(" %s══════════════════════════════════════════════════════════════════════════%s\n\n", ansiDim, ansiReset))

	b.WriteString(fmt.Sprintf(" %s%s[ TOP DISPATCHES ]%s\n\n", ansiGreen, ansiBold, ansiReset))
	count := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		tags := ""
		if len(p.Tags) > 0 {
			tags = fmt.Sprintf(" %s[%s]%s", ansiDim, strings.ToUpper(strings.Join(p.Tags, " ")), ansiReset)
		}
		b.WriteString(fmt.Sprintf("  %s%d.%s %s%s%s%s\n", ansiYellow, count+1, ansiReset, ansiBold, p.Title, ansiReset, tags))
		if p.Description != "" {
			b.WriteString(fmt.Sprintf("     %s%s%s\n", ansiDim, p.Description, ansiReset))
		}
		b.WriteString(fmt.Sprintf("    curl -sL %s/p/%s\n\n", r.Host, p.Slug))
		count++
		if count >= 8 {
			break
		}
	}

	b.WriteString(fmt.Sprintf(" %s%s[ TERMINAL ENDPOINTS ]%s\n\n", ansiMagenta, ansiBold, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/daily          %s# Stream full daily tech briefing%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/recipes        %s# Production Linux & sysadmin recipes%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/go           %s# Stream Go backend dispatches%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/linux        %s# Stream Linux kernel & OS stories%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/t/python       %s# Stream Python & backend articles%s\n", r.Host, ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("  curl -sL %s/p/<slug>       %s# Read full article in terminal%s\n\n", r.Host, ansiDim, ansiReset))

	b.WriteString(fmt.Sprintf(" %sTip: Run 'curl -sL %s/daily | less -R' for paginated reading.%s\n\n", ansiDim, r.Host, ansiReset))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

func (h *Handler) CLITag(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("nocolor") != "1"
	tag := strings.ToLower(chi.URLParam(r, "tag"))
	if tag == "" {
		tag = strings.ToLower(r.URL.Query().Get("tag"))
	}
	posts := h.AllPosts()

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s%s[ DAEMONTALK TAG STREAM: #%s ]%s\n\n", ansiCyan, ansiBold, tag, ansiReset))

	matched := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		hasTag := false
		for _, t := range p.Tags {
			if strings.EqualFold(t, tag) {
				hasTag = true
				break
			}
		}
		if hasTag {
			b.WriteString(fmt.Sprintf("  %s•%s %s%s%s (%s)\n", ansiGreen, ansiReset, ansiBold, p.Title, ansiReset, p.Date.Format("2006-01-02")))
			if p.Description != "" {
				b.WriteString(fmt.Sprintf("    %s%s%s\n", ansiDim, p.Description, ansiReset))
			}
			b.WriteString(fmt.Sprintf("    %scurl -sL %s/p/%s%s\n\n", ansiYellow, r.Host, p.Slug, ansiReset))
			matched++
		}
	}

	if matched == 0 {
		b.WriteString(fmt.Sprintf("  No articles found matching tag '#%s'.\n", tag))
		b.WriteString(fmt.Sprintf("  Available tags: linux, rust, go, python, devops, docker, security, storage, tools\n\n"))
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

func (h *Handler) CLIPost(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("nocolor") != "1"
	slug := chi.URLParam(r, "slug")
	slug = strings.TrimSuffix(slug, ".txt")

	p, ok := post.FindBySlug(h.AllPosts(), slug)
	if !ok {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("Post not found: %s\n", slug)))
		return
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s%s# %s%s\n\n", ansiCyan, ansiBold, p.Title, ansiReset))
	b.WriteString(fmt.Sprintf("%sDate:%s %s | %sLang:%s %s | %sTags:%s [%s]\n", ansiDim, ansiReset, p.Date.Format("2006-01-02"), ansiDim, ansiReset, p.Lang, ansiDim, ansiReset, strings.Join(p.Tags, ", ")))
	b.WriteString(fmt.Sprintf("%s══════════════════════════════════════════════════════════════════════════%s\n\n", ansiDim, ansiReset))

	if p.Description != "" {
		b.WriteString(fmt.Sprintf("%s> %s%s\n\n", ansiYellow, p.Description, ansiReset))
	}

	// Render raw markdown/html body
	b.WriteString(string(p.Body))
	b.WriteString(fmt.Sprintf("\n\n%s══════════════════════════════════════════════════════════════════════════%s\n", ansiDim, ansiReset))
	b.WriteString(fmt.Sprintf("%sRead on web: https://%s/blog/%s%s\n\n", ansiDim, r.Host, p.Slug, ansiReset))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

func (h *Handler) CLIRecipes(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("nocolor") != "1"

	recipes := `
 [ PRODUCTION-READY LINUX & eBPF RECIPES ]
 ══════════════════════════════════════════════════════════════════════════

 1. Trace Top System Calls by Process in Real-Time (bpftrace)
    sudo bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count(); }'

 2. Monitor Files Opened Across the Entire System (bpftrace)
    sudo bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf("%s opened %s\n", comm, str(args->filename)); }'

 3. Find Out Which Process is Creating Heavy Disk Write I/O (iotop / pidstat)
    sudo pidstat -d 1 5
    sudo iotop -oPa

 4. Trace TCP Packet Drops in Kernel (pwru / dropwatch / bpftrace)
    sudo bpftrace -e 'kprobe:kfree_skb { @[kstack] = count(); }'

 5. Find Large Deleted Files Still Held Open by Processes (Reclaim Disk Space)
    sudo lsof +L1 | grep deleted

 6. Check Inode Saturation when 'No space left on device' (even with 50% free disk)
    df -i
    sudo find /var -xdev -printf '%h\n' | sort | uniq -c | sort -k1 -nr | head -n 10

 7. Profile Go / Rust Application CPU on Linux with perf & FlameGraph
    sudo perf record -F 99 -p <PID> -g -- sleep 30
    sudo perf script | stackcollapse-perf.pl | flamegraph.pl > flame.svg

 8. Find Socket Listening Processes and Zero-Latency Network Connections
    ss -tulpn
    ss -s
`

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s%s%s%s\n", ansiGreen, ansiBold, recipes, ansiReset))
	b.WriteString(fmt.Sprintf(" %sMore recipes and interactive challenges at: https://%s/terminal%s\n\n", ansiCyan, r.Host, ansiReset))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

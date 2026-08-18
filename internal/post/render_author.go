package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

func renderAuthorHTML(rawContent string) string {
	meta := make(map[string]string)
	lines := strings.Split(rawContent, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			meta[k] = v
		}
	}

	name := meta["name"]
	if name == "" {
		name = "Dafa Gareth"
	}
	role := meta["role"]
	if role == "" {
		role = "Software Engineer"
	}
	avatar := meta["avatar"]
	if avatar == "" || avatar == "/static/logo/logo-dark.png" || avatar == "/static/logo/logo-light.png" || avatar == "/static/logo/logo-dark.webp" || avatar == "/static/logo/logo-light.webp" {
		avatar = "/static/logo/icon-dark.png"
	}
	bio := meta["bio"]
	if bio == "" {
		bio = "Software Engineer yang berfokus pada sistem terdistribusi, rekayasa kernel Linux, dan optimasi performa backend Go/Rust."
	}

	var buf bytes.Buffer
	buf.WriteString("\n<div class=\"post-author-card my-8 not-prose p-4 sm:p-5 border border-border bg-surface rounded-none flex flex-col gap-3.5\">\n")
	buf.WriteString(fmt.Sprintf(`  <!-- Top row: Avatar (PFP) on left, Name & Role directly beside it -->
  <div class="flex items-center gap-3.5 sm:gap-4">
    <div class="w-12 h-12 sm:w-14 sm:h-14 shrink-0 rounded-full border border-border bg-chip/40 overflow-hidden flex items-center justify-center p-1 shadow-sm">
      <img src="%s" alt="%s" class="w-full h-full object-contain block" loading="lazy" />
    </div>
    <div class="flex flex-col min-w-0">
      <h3 class="author-name text-[1.05em] sm:text-[1.1em] font-bold text-text uppercase tracking-wider leading-tight truncate">%s</h3>
      <span class="author-role text-[0.8em] sm:text-[0.82em] font-mono text-muted mt-0.5 truncate">%s</span>
    </div>
  </div>
  <!-- Bottom row: Description spanning horizontally underneath -->
  <p class="author-bio text-[0.88em] sm:text-[0.92em] text-muted leading-relaxed">%s</p>
  <div class="flex flex-wrap items-center gap-2 text-[0.8em] font-mono pt-2 border-t border-border/40">
`, html.EscapeString(avatar), html.EscapeString(name), html.EscapeString(name), html.EscapeString(role), html.EscapeString(bio)))

	socialKeys := []string{"github", "gh", "x", "twitter", "linkedin", "li", "email", "mail", "website", "site", "blog", "youtube", "yt", "instagram", "ig", "threads", "bluesky", "bsky", "telegram", "tg", "gitlab", "discord"}
	seen := make(map[string]bool)

	for _, k := range socialKeys {
		val, exists := meta[k]
		if !exists || val == "" {
			continue
		}
		u, label, icon := getSocialLink(k, val)
		if u == "" || seen[label] {
			continue
		}
		seen[label] = true
		isExt := !strings.HasPrefix(u, "mailto:")
		relAttr := ""
		targetAttr := ""
		if isExt {
			targetAttr = ` target="_blank"`
			relAttr = ` rel="noopener noreferrer"`
		}
		buf.WriteString(fmt.Sprintf(`    <a href="%s"%s%s title="%s" aria-label="%s" class="w-7 h-7 inline-flex items-center justify-center border border-border bg-chip/40 text-muted hover:text-text hover:bg-hover hover:border-text transition-colors">
      %s
    </a>
`, html.EscapeString(u), targetAttr, relAttr, html.EscapeString(label), html.EscapeString(label), icon))
	}

	buf.WriteString(`  </div>
</div>
`)
	return buf.String()
}

func getSocialLink(key, val string) (url, label, svgIcon string) {
	val = strings.TrimSpace(val)
	if val == "" {
		return "", "", ""
	}
	cleanUsername := strings.TrimPrefix(val, "@")

	switch key {
	case "github", "gh":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://github.com/" + cleanUsername
		}
		label = "GitHub"
		svgIcon = GetIcon(IconGitHub, "w-3.5 h-3.5 shrink-0")

	case "x", "twitter":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://x.com/" + cleanUsername
		}
		label = "X"
		svgIcon = GetIcon(IconTwitter, "w-3.5 h-3.5 shrink-0")

	case "email", "mail":
		if strings.HasPrefix(val, "mailto:") {
			url = val
			label = strings.TrimPrefix(val, "mailto:")
		} else {
			url = "mailto:" + val
			label = val
		}
		svgIcon = GetIcon(IconEmail, "w-3.5 h-3.5 shrink-0")

	case "website", "site", "blog", "url", "link":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://" + val
		}
		label = "Website"
		svgIcon = GetIcon(IconWebsite, "w-3.5 h-3.5 shrink-0")

	case "linkedin", "li":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://linkedin.com/in/" + cleanUsername
		}
		label = "LinkedIn"
		svgIcon = GetIcon(IconLinkedIn, "w-3.5 h-3.5 shrink-0")

	case "youtube", "yt":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://youtube.com/@" + cleanUsername
		}
		label = "YouTube"
		svgIcon = GetIcon(IconYouTube, "w-3.5 h-3.5 shrink-0")

	case "instagram", "ig":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://instagram.com/" + cleanUsername
		}
		label = "Instagram"
		svgIcon = GetIcon(IconInstagram, "w-3.5 h-3.5 shrink-0")

	case "threads":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://threads.net/@" + cleanUsername
		}
		label = "Threads"
		svgIcon = GetIcon(IconThreads, "w-3.5 h-3.5 shrink-0")

	case "bluesky", "bsky":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://bsky.app/profile/" + cleanUsername
		}
		label = "Bluesky"
		svgIcon = GetIcon(IconBluesky, "w-3.5 h-3.5 shrink-0")

	case "telegram", "tg":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://t.me/" + cleanUsername
		}
		label = "Telegram"
		svgIcon = GetIcon(IconTelegram, "w-3.5 h-3.5 shrink-0")

	case "gitlab":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://gitlab.com/" + cleanUsername
		}
		label = "GitLab"
		svgIcon = GetIcon(IconGitLab, "w-3.5 h-3.5 shrink-0")

	case "discord":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://discord.gg/" + val
		}
		label = "Discord"
		svgIcon = GetIcon(IconDiscord, "w-3.5 h-3.5 shrink-0")
	}
	return url, label, svgIcon
}

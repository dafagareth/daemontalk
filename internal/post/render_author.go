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
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>`

	case "x", "twitter":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://x.com/" + cleanUsername
		}
		label = "X"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.744l7.737-8.835L1.254 2.25H8.08l4.253 5.622 5.911-5.622zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>`

	case "email", "mail":
		if strings.HasPrefix(val, "mailto:") {
			url = val
			label = strings.TrimPrefix(val, "mailto:")
		} else {
			url = "mailto:" + val
			label = val
		}
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"></rect><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"></path></svg>`

	case "website", "site", "blog", "url", "link":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://" + val
		}
		label = "Website"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>`

	case "linkedin", "li":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://linkedin.com/in/" + cleanUsername
		}
		label = "LinkedIn"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.779-1.75-1.75s.784-1.75 1.75-1.75 1.75.779 1.75 1.75-.784 1.75-1.75 1.75zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z"/></svg>`

	case "youtube", "yt":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://youtube.com/@" + cleanUsername
		}
		label = "YouTube"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>`

	case "instagram", "ig":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://instagram.com/" + cleanUsername
		}
		label = "Instagram"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="20" x="2" y="2" rx="5" ry="5"></rect><path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"></path><line x1="17.5" x2="17.51" y1="6.5" y2="6.5"></line></svg>`

	case "threads":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://threads.net/@" + cleanUsername
		}
		label = "Threads"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12.786 2.003c-5.918-.28-10.783 4.29-10.783 10.098 0 5.617 4.544 10.021 10.098 10.021 4.887 0 8.847-3.415 9.771-8.084h-2.274c-.846 3.42-3.83 5.922-7.497 5.922-4.328 0-7.857-3.432-7.857-7.859 0-4.498 3.662-8.031 8.243-7.859 4.053.153 7.302 3.415 7.424 7.469h2.245c-.125-5.285-4.372-9.528-9.37-9.708zm2.638 8.02c-.347-.367-.847-.577-1.393-.585-.972-.015-1.875.568-2.247 1.45-.16.38-.247.788-.258 1.212.525-.084 1.063-.125 1.613-.125 1.123 0 1.99.42 2.384 1.15.318.59.318 1.315 0 1.905-.38.71-1.155 1.1-2.18 1.1-.95 0-1.793-.453-2.193-1.182a3.94 3.94 0 0 1-.355-1.664v-1.574c0-1.513.648-2.873 1.777-3.732 1.052-.8 2.428-1.2 3.982-1.156 1.89.054 3.417.822 4.298 2.16.76 1.155 1.024 2.707.762 4.49-.336 2.29-1.706 4.144-3.756 5.09-1.355.626-2.902.903-4.473.803-2.536-.16-4.652-1.563-5.524-3.658l2.105-1.226c.535 1.274 1.83 2.13 3.346 2.226 1.036.066 2.053-.117 2.948-.529 1.342-.62 2.246-1.864 2.48-3.414.2-1.305.009-2.413-.53-3.208-.565-.833-1.533-1.298-2.72-1.31zm-3.318 4.301c.025.447.18.812.462 1.09.275.27.653.413 1.094.413.494 0 .868-.164 1.051-.463.161-.262.167-.61.015-.924-.198-.407-.702-.697-1.421-.697-.417 0-.82.032-1.201.095v.486z"/></svg>`

	case "bluesky", "bsky":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://bsky.app/profile/" + cleanUsername
		}
		label = "Bluesky"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 10.8c-1.087-2.114-4.046-6.053-6.798-7.995C2.566 1.01 1.05 1.638.384 3.014c-.663 1.374-.183 4.225.86 6.554 1.042 2.329 2.766 4.308 4.675 5.568-1.909 1.26-3.633 3.239-4.675 5.568-1.043 2.329-1.523 5.18-.86 6.554.666 1.376 2.182 2.004 4.818.209C7.954 25.52 10.913 21.58 12 19.467c1.087 2.114 4.046 6.053 6.798 7.995 2.636 1.795 4.152 1.167 4.818-.209.663-1.374.183-4.225-.86-6.554-1.042-2.329-2.766-4.308-4.675-5.568 1.909-1.26 3.633-3.239 4.675-5.568 1.043-2.329 1.523-5.18.86-6.554-.666-1.376-2.182-2.004-4.818-.209C16.046 4.747 13.087 8.686 12 10.8z"/></svg>`

	case "telegram", "tg":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://t.me/" + cleanUsername
		}
		label = "Telegram"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 0 0-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.08-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .38z"/></svg>`

	case "gitlab":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://gitlab.com/" + cleanUsername
		}
		label = "GitLab"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="m23.6 9.59-1.22-3.76a.86.86 0 0 0-1.63 0l-1.22 3.76H4.47L3.25 5.83a.86.86 0 0 0-1.63 0L.4 9.59a1.73 1.73 0 0 0 .63 1.93l10.42 7.57a1 1 0 0 0 1.1 0L23 11.52a1.73 1.73 0 0 0 .6-1.93z"/></svg>`

	case "discord":
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			url = val
		} else {
			url = "https://discord.gg/" + val
		}
		label = "Discord"
		svgIcon = `<svg class="w-3.5 h-3.5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994.021-.041.001-.09-.041-.106a13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.929 1.793 8.18 1.793 12.061 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.893.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.028zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/></svg>`
	}
	return url, label, svgIcon
}

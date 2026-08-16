package post

import (
	"regexp"
	"strings"
)

var (
	reCarouselFenced   = regexp.MustCompile("(?s)```(?:carousel|slider)\\s*\\n(.*?)\\n```")
	reGalleryFenced    = regexp.MustCompile("(?s)```gallery\\s*\\n(.*?)\\n```")
	reFAQFenced        = regexp.MustCompile("(?s)```faq\\s*\\n(.*?)\\n```")
	reAuthorFenced     = regexp.MustCompile("(?s)```author\\s*\\n(.*?)\\n```")
	reReferencesFenced = regexp.MustCompile("(?s)```(?:references|refs|ref)\\s*\\n(.*?)\\n```")
	reCalloutFenced    = regexp.MustCompile("(?s)```callout(?:\\s+([a-zA-Z]+))?\\s*\\n(.*?)\\n```")
	reGitHubAlert      = regexp.MustCompile(`(?m)^>\s*\[\!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\][^\n]*\n((?:>[^\n]*(?:\n|$))+)`)
	reTabsFenced       = regexp.MustCompile("(?s)```(?:tabs|files|code-tabs)\\s*\\n(.*?)\\n```")
	reLinkFenced       = regexp.MustCompile("(?s)```(?:link|bookmark|card)\\s*\\n(.*?)\\n```")
	reStatFenced       = regexp.MustCompile("(?s)```(?:stat|stats|metrics)\\s*\\n(.*?)\\n```")
	reMarkdownImage    = regexp.MustCompile(`!\[([^\]]*)\]\(([^)"\s]+)(?:\s+(?:"([^"]*)"|'([^']*)'))?\)`)
)

// preprocessMarkdown handles custom media blocks such as ```carousel, ```gallery, ```faq, ```author, ```references, ```callout, ```tabs, ```link, and ```stat.
func preprocessMarkdown(src []byte) []byte {
	srcStr := string(src)

	// Replace GitHub style alert blockquotes: > [!NOTE] ...
	srcStr = reGitHubAlert.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reGitHubAlert.FindStringSubmatch(m)
		if len(sub) < 3 {
			return m
		}
		alertType := sub[1]
		bodyLines := strings.Split(sub[2], "\n")
		var cleanLines []string
		for _, l := range bodyLines {
			l = strings.TrimPrefix(l, ">")
			l = strings.TrimPrefix(l, " ")
			cleanLines = append(cleanLines, l)
		}
		return renderCalloutHTML(alertType, strings.Join(cleanLines, "\n"))
	})

	// Replace ```callout ... ```
	srcStr = reCalloutFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reCalloutFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		rawContent := sub[len(sub)-1]
		alertType := "NOTE"
		if len(sub) >= 3 && sub[1] != "" {
			alertType = sub[1]
		}
		lines := strings.Split(rawContent, "\n")
		var bodyLines []string
		for _, l := range lines {
			trimmed := strings.TrimSpace(l)
			if strings.HasPrefix(strings.ToLower(trimmed), "type:") {
				alertType = strings.TrimSpace(trimmed[5:])
			} else {
				bodyLines = append(bodyLines, l)
			}
		}
		return renderCalloutHTML(alertType, strings.Join(bodyLines, "\n"))
	})

	// Replace ```references ... ```
	srcStr = reReferencesFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reReferencesFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderReferencesHTML(sub[1])
	})

	// Replace ```tabs ... ```
	srcStr = reTabsFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reTabsFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderTabsHTML(sub[1])
	})

	// Replace ```link ... ```
	srcStr = reLinkFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reLinkFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderLinkHTML(sub[1])
	})

	// Replace ```stat ... ```
	srcStr = reStatFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reStatFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderStatHTML(sub[1])
	})

	// Replace ```carousel ... ```
	srcStr = reCarouselFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reCarouselFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		items := extractMediaItems(sub[1])
		if len(items) == 0 {
			return m
		}
		return renderCarouselHTML(items)
	})

	// Replace ```gallery ... ```
	srcStr = reGalleryFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reGalleryFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		items := extractMediaItems(sub[1])
		if len(items) == 0 {
			return m
		}
		return renderGalleryHTML(items)
	})

	// Replace ```faq ... ```
	srcStr = reFAQFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reFAQFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderFAQHTML(sub[1])
	})

	// Replace ```author ... ```
	srcStr = reAuthorFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reAuthorFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderAuthorHTML(sub[1])
	})

	return []byte(srcStr)
}

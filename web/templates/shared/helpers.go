package shared

import (
	"fmt"
	"strings"
)

func Prefix(lang string) string {
	return ""
}

func PageTitle(page string) string {
	suffix := "daemontalk"
	switch page {
	case "home", "":
		return suffix
	case "projects":
		return "Projects · " + suffix
	case "blog":
		return "Blog · " + suffix
	case "about":
		return "About · " + suffix
	case "colophon":
		return "Colophon · " + suffix
	case "uses":
		return "Uses · " + suffix
	case "now":
		return "Now · " + suffix
	case "saved":
		return "Reading List · " + suffix
	case "stats":
		return "Stats · " + suffix
	case "guestbook":
		return "Guestbook · " + suffix
	case "search":
		return "Search · " + suffix
	case "links":
		return "Links · " + suffix
	case "404":
		return "404 · Page Not Found · " + suffix
	case "403":
		return "403 · Access Forbidden · " + suffix
	case "500":
		return "500 · Internal Server Error · " + suffix
	default:
		if strings.HasPrefix(page, "#") {
			return page + " · " + suffix
		}
		return page
	}
}

func HreflangEN(path string) string {
	p := strings.TrimPrefix(path, "/id")
	if p == "" {
		return "/"
	}
	return p
}

func IsSocketApp(path string) bool {
	p := strings.TrimPrefix(path, "/id")
	return p == "/socket" || (strings.HasPrefix(p, "/socket/") && p != "/socket/new")
}

func CoverSourceName(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	lower := strings.ToLower(rawURL)
	switch {
	case strings.Contains(lower, "unsplash.com"):
		return "Unsplash"
	case strings.Contains(lower, "github.com"):
		return "GitHub"
	case strings.Contains(lower, "wikimedia.org") || strings.Contains(lower, "wikipedia.org"):
		return "Wikimedia"
	case strings.Contains(lower, "pexels.com"):
		return "Pexels"
	}
	s := strings.TrimPrefix(rawURL, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "www.")
	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}
	if s != "" {
		return s
	}
	return "Source"
}

func BibtexFromMeta(slug, title, date string) string {
	year := "2026"
	if len(date) >= 4 {
		year = date[:4]
	}
	return fmt.Sprintf(`@article{daemontalk_%s,
  title = "{%s}",
  author = {Daemontalk Team},
  year = "%s",
  url = "https://daemontalk.com/blog/%s"
}`, strings.ReplaceAll(slug, "-", "_"), title, year, slug)
}

func DemoLabel(url string) string {
	if strings.Contains(url, "/releases") {
		return "Releases"
	}
	return "Demo"
}

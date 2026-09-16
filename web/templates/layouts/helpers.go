package layouts

import (
	"strings"

	"daemontalk/web/templates/shared"
)

func isHomeOrTagFilter(path string) bool {
	p := strings.TrimPrefix(path, "/id")
	return p == "" || p == "/" || p == "/blog" || p == "/blog/" || strings.HasPrefix(p, "/blog/tag/") || p == "/search" || p == "/search/"
}

func isSocketPath(path string) bool {
	p := strings.TrimPrefix(path, "/id")
	return p == "/socket" || strings.HasPrefix(p, "/socket/")
}

func isBlogPostPath(path string) bool {
	p := strings.TrimPrefix(path, "/id")
	return strings.HasPrefix(p, "/blog/") && !strings.HasPrefix(p, "/blog/tag/") && p != "/blog" && p != "/blog/" && p != "/blog/posts"
}

func isMainHome(path string) bool {
	p := strings.TrimPrefix(path, "/id")
	return p == "" || p == "/" || p == "/blog" || p == "/blog/"
}

func postTitleFromPage(page string) string {
	return strings.TrimSuffix(page, " · daemontalk")
}

func prefix(lang string) string {
	return shared.Prefix(lang)
}

func pageTitle(page string) string {
	return shared.PageTitle(page)
}

func hreflangEN(path string) string {
	return shared.HreflangEN(path)
}

func isSocketApp(path string) bool {
	return shared.IsSocketApp(path)
}

func assetURL(path string) string {
	return shared.AssetURL(path)
}

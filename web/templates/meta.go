package templates

// SiteBaseURL is the base canonical domain for social/OG tags and links.
// It defaults to production and can be configured at startup via BASE_URL.
var SiteBaseURL = "https://www.daemontalk.com"

// AssetVersion is a cache-busting token appended to static CSS/JS URLs. It is
// set once at startup (see main.go) so a rebuild of main.css is picked up by
// browsers immediately instead of serving a stale cached copy.
var AssetVersion = "dev"

func assetURL(path string) string {
	return path + "?v=" + AssetVersion
}

type PageMeta struct {
	Description   string
	Image         string
	Type          string // "website" or "article"
	JSONLD        string // optional schema.org JSON-LD, injected raw into <head>
	PublishedTime string // RFC3339, for article:published_time
	Author        string
}

func (m PageMeta) ogType() string {
	if m.Type != "" {
		return m.Type
	}
	return "website"
}

func (m PageMeta) ogDesc() string {
	if m.Description != "" {
		return m.Description
	}
	return "An independent technology publication and research notebook exploring modern computing, software architecture, and systems engineering."
}

func (m PageMeta) ogImage() string {
	if m.Image != "" {
		return m.Image
	}
	return SiteBaseURL + "/og.png"
}

// AbsoluteURL turns a site-relative path ("/static/...") into a full URL for
// use in social/OG tags. Already-absolute URLs are returned unchanged.
func AbsoluteURL(path string) string {
	if path == "" {
		return ""
	}
	if len(path) >= 4 && (path[:4] == "http") {
		return path
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return SiteBaseURL + path
}

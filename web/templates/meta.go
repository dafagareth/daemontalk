package templates

const siteBaseURL = "https://www.daemontalk.com"

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
	return "Open engineering notebook and learning log. Exploring Go backend, Python, Linux systems, and distributed architecture."
}

func (m PageMeta) ogImage() string {
	if m.Image != "" {
		return m.Image
	}
	return siteBaseURL + "/og.png"
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
	return siteBaseURL + path
}

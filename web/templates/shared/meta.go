package shared

import (
	"fmt"
	"os"
	"strings"
)

var SiteBaseURL = "https://www.daemontalk.com"

var AssetVersion = "dev"

func AssetURL(path string) string {
	cleanPath := strings.TrimPrefix(path, "/")
	if fi, err := os.Stat("web/" + cleanPath); err == nil {
		return path + fmt.Sprintf("?v=%d", fi.ModTime().Unix())
	}
	return path + "?v=" + AssetVersion
}

type PageMeta struct {
	Description   string
	Image         string
	URL           string
	Type          string
	JSONLD        string
	PublishedTime string
	Author        string
}

func (m PageMeta) OgURL(currentPath string) string {
	if m.URL != "" {
		return m.URL
	}
	return SiteBaseURL + currentPath
}

func (m PageMeta) OgType() string {
	if m.Type != "" {
		return m.Type
	}
	return "website"
}

func (m PageMeta) OgDesc() string {
	if m.Description != "" {
		return m.Description
	}
	return "An independent tech publication and learning space exploring modern software, systems, and the digital world for everyone."
}

func (m PageMeta) OgImage() string {
	if m.Image != "" {
		return m.Image
	}
	return SiteBaseURL + "/og.png"
}

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

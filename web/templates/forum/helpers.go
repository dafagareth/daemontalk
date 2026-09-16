package forum

import (
	"fmt"
	"net/url"
	"time"
)

func prefix(_ string) string {
	return ""
}

func timeAgoOrDate(d time.Time, _ string) string {
	dur := time.Since(d)
	if dur < 0 {
		dur = 0
	}
	mins := int(dur.Minutes())
	if mins < 1 {
		return "just now"
	}
	if mins < 60 {
		return fmt.Sprintf("%dm ago", mins)
	}
	hours := int(dur.Hours())
	if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}
	days := hours / 24
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}
	return d.Format("02 Jan 2006")
}

func getFilterURL(prefixStr, tag, sort string) string {
	base := prefixStr + "/socket"
	v := url.Values{}
	if tag != "" {
		v.Set("tag", tag)
	}
	if sort != "" && sort != "latest" {
		v.Set("sort", sort)
	}
	if len(v) > 0 {
		return base + "?" + v.Encode()
	}
	return base
}

func getPageURL(prefixStr, tag, sort string, p int) string {
	base := prefixStr + "/socket"
	v := url.Values{}
	if tag != "" {
		v.Set("tag", tag)
	}
	if sort != "" && sort != "latest" {
		v.Set("sort", sort)
	}
	if p > 1 {
		v.Set("page", fmt.Sprintf("%d", p))
	}
	if len(v) > 0 {
		return base + "?" + v.Encode()
	}
	return base
}

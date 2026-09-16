package shared

import (
	"fmt"
	"time"

	"daemontalk/internal/post"
)

func FmtDate(p post.Post, lang string) string {
	return FmtDateFromTime(p.Date, lang)
}

func FmtDateFromTime(d time.Time, _ string) string {
	return d.Format("02 January 2006") + " WIB"
}

func FmtShortDate(d time.Time, _ string) string {
	return d.Format("02 Jan 2006")
}

func FmtDayMonth(d time.Time, _ string) string {
	return d.Format("02 Jan")
}

func FmtDateTime(d time.Time, _ string) string {
	return d.Format("02 Jan 2006, 15:04") + " WIB"
}

func FmtFullDate(d time.Time, _ string) string {
	return d.Format("Monday, 02 January 2006") + " WIB"
}

func TimeAgoOrDate(d time.Time, lang string) string {
	dur := time.Since(d)
	if dur < 0 {
		dur = 0
	}
	mins := int(dur.Minutes())
	if mins < 1 {
		return "just now"
	}
	if mins < 60 {
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	}
	hrs := int(dur.Hours())
	if hrs < 24 {
		if hrs == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hrs)
	}
	days := int(dur.Hours() / 24)
	if days < 30 {
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	return FmtShortDate(d, lang)
}

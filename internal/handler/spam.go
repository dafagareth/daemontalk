package handler

import (
	"strings"
	"unicode"
)

const spamThreshold = 30

var spamKeywords = []string{
	"casino", "buy now", "click here", "free money",
	"earn money", "work from home", "100% free",
	"make money fast", "lose weight fast", "diet pill",
	"follow me", "subscribe now", "limited offer", "gacor",
}

// spamScore returns a heuristic spam likelihood score for a comment.
// Score > spamThreshold should be silently rejected.
func spamScore(name, body string) int {
	score := 0

	// Multiple HTTP links in body
	score += strings.Count(body, "http") * 15

	// Very short name (< 2 non-space chars)
	if len(strings.TrimSpace(name)) < 2 {
		score += 20
	}

	// Spam keywords
	lowerBody := strings.ToLower(body)
	for _, kw := range spamKeywords {
		if strings.Contains(lowerBody, kw) {
			score += 10
		}
	}

	// High ALL-CAPS ratio (> 60% of letters are uppercase)
	var letters, caps int
	for _, r := range body {
		if unicode.IsLetter(r) {
			letters++
			if unicode.IsUpper(r) {
				caps++
			}
		}
	}
	if letters > 15 && caps*100/letters > 60 {
		score += 20
	}

	return score
}

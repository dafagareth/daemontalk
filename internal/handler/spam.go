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

func spamScore(name, body string) int {
	score := 0

	score += strings.Count(body, "http") * 15

	if len(strings.TrimSpace(name)) < 2 {
		score += 20
	}

	lowerBody := strings.ToLower(body)
	for _, kw := range spamKeywords {
		if strings.Contains(lowerBody, kw) {
			score += 10
		}
	}

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

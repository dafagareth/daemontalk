package handler

import "daemontalk/internal/handler/common"

const spamThreshold = common.SpamThreshold

var spamKeywords = common.SpamKeywords

func spamScore(name, body string) int {
	return common.SpamScore(name, body)
}

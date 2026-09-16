package handler

import "daemontalk/internal/handler/common"

const (
	CookieAdminToken         = common.CookieAdminToken
	CookieVisitorID          = common.CookieVisitorID
	CookieReactedPrefix      = common.CookieReactedPrefix
	CookieViewCooldownPrefix = common.CookieViewCooldownPrefix

	CookieAdminMaxAge        = common.CookieAdminMaxAge
	CookieReactionMaxAge     = common.CookieReactionMaxAge
	CookieViewCooldownMaxAge = common.CookieViewCooldownMaxAge
	CookieVisitorExpiryYears = common.CookieVisitorExpiryYears

	DefaultPostsPerPage   = common.DefaultPostsPerPage
	MaxMarkdownUploadSize = common.MaxMarkdownUploadSize
	MaxImageUploadSize    = common.MaxImageUploadSize

	GuestbookSlug = common.GuestbookSlug
)

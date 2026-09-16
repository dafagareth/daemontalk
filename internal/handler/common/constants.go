package common

const (
	CookieAdminToken         = "admin_token"
	CookieVisitorID          = "visitor_id"
	CookieReactedPrefix      = "reacted_"
	CookieViewCooldownPrefix = "v_post_"

	CookieAdminMaxAge        = 60 * 60 * 24 * 30
	CookieReactionMaxAge     = 86400 * 365
	CookieViewCooldownMaxAge = 3600 * 12
	CookieVisitorExpiryYears = 10

	DefaultPostsPerPage   = 14
	MaxMarkdownUploadSize = 20 << 20
	MaxImageUploadSize    = 10 << 20

	GuestbookSlug = "__guestbook__"
)

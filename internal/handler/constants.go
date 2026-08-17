package handler

const (
	// Cookie names
	CookieAdminToken         = "admin_token"
	CookieVisitorID          = "visitor_id"
	CookieReactedPrefix      = "reacted_"
	CookieViewCooldownPrefix = "v_post_"

	// Cookie MaxAge durations in seconds
	CookieAdminMaxAge        = 60 * 60 * 24 * 30 // 30 days
	CookieReactionMaxAge     = 86400 * 365       // 1 year
	CookieViewCooldownMaxAge = 3600 * 12         // 12 hours
	CookieVisitorExpiryYears = 10                // 10 years

	// Pagination & Limits
	DefaultPostsPerPage   = 14
	MaxMarkdownUploadSize = 20 << 20 // 20 MB
	MaxImageUploadSize    = 10 << 20 // 10 MB

	// System slugs & identifiers
	GuestbookSlug = "__guestbook__"
)

package auth

import "time"

// User represents an authenticated member on Daemontalk.
type User struct {
	ID          int64  `json:"id"`
	Provider    string `json:"provider"`
	ProviderID  string `json:"provider_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	GitHubURL   string `json:"github_url"`
	// "member", "moderator", "admin"
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session represents a persistent login session backed by an encrypted cookie token.
type Session struct {
	TokenHash string    `json:"token_hash"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

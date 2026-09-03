package comment

import "time"

// Comment represents a visitor comment or reply.
type Comment struct {
	ID         int64     `json:"id"`
	PostSlug   string    `json:"post_slug"`
	Name       string    `json:"name"`
	Body       string    `json:"body"`
	ParentID   *int64    `json:"parent_id,omitempty"`
	UserID     *int64    `json:"user_id,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	IsVerified bool      `json:"is_verified"`
	GitHubURL  string    `json:"github_url,omitempty"`
	IsReported bool      `json:"is_reported"`
	CreatedAt  time.Time `json:"created_at"`
	Replies    []Comment `json:"replies,omitempty"`
}

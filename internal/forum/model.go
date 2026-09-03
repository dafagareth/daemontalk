package forum

import (
	"html/template"
	"time"
)

// Category defines available technical discussion categories.
type Category struct {
	Slug        string
	NameID      string
	NameEN      string
	Description string
	Icon        string
}

var Categories = []Category{
	{Slug: "all", NameID: "Semua Topik", NameEN: "All Topics", Description: "Semua diskusi teknologi & sistem"},
	{Slug: "qna", NameID: "Tanya Jawab (Q&A)", NameEN: "Q&A / Debugging", Description: "Pertanyaan teknis & pemecahan masalah"},
	{Slug: "architecture", NameID: "Arsitektur Sistem", NameEN: "Systems Architecture", Description: "Desain sistem terdistribusi & backend"},
	{Slug: "kernel", NameID: "Linux & Kernel", NameEN: "Linux & Kernel", Description: "Eksplorasi low-level, eBPF & OS internals"},
	{Slug: "go", NameID: "Go & Concurrency", NameEN: "Go & Concurrency", Description: "Runtime Go, memory model & performa"},
	{Slug: "incident", NameID: "Kasus Insiden (RCA)", NameEN: "Incident & RCA", Description: "Post-mortem & investigasi kegagalan sistem"},
}

// Topic represents a discussion thread or Q&A question in Daemontalk.
type Topic struct {
	ID             int64         `json:"id"`
	UserID         int64         `json:"user_id"`
	AuthorName     string        `json:"author_name"`
	AuthorUsername string        `json:"author_username"`
	AuthorAvatar   string        `json:"author_avatar"`
	AuthorGitHub   string        `json:"author_github"`
	Title          string        `json:"title"`
	Slug           string        `json:"slug"`
	Category       string        `json:"category"`
	Tags           []string      `json:"tags"`
	BodyMD         string        `json:"body_md"`
	BodyHTML       template.HTML `json:"body_html"`
	SolvedReplyID  int64         `json:"solved_reply_id,omitempty"`
	IsSolved       bool          `json:"is_solved"`
	ViewsCount     int           `json:"views_count"`
	VotesCount     int           `json:"votes_count"`
	RepliesCount   int           `json:"replies_count"`
	Pinned         bool          `json:"pinned"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	UserVoted      bool          `json:"user_voted,omitempty"`
	IsOwner        bool          `json:"is_owner,omitempty"`
}

// Reply represents a reply or potential solution to a discussion topic.
type Reply struct {
	ID             int64         `json:"id"`
	TopicID        int64         `json:"topic_id"`
	ParentID       int64         `json:"parent_id,omitempty"`
	UserID         int64         `json:"user_id"`
	AuthorName     string        `json:"author_name"`
	AuthorUsername string        `json:"author_username"`
	AuthorAvatar   string        `json:"author_avatar"`
	AuthorGitHub   string        `json:"author_github"`
	BodyMD         string        `json:"body_md"`
	BodyHTML       template.HTML `json:"body_html"`
	IsSolution     bool          `json:"is_solution"`
	VotesCount     int           `json:"votes_count"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	UserVoted      bool          `json:"user_voted,omitempty"`
	IsOwner        bool          `json:"is_owner,omitempty"`
	Children       []*Reply      `json:"children,omitempty"`
}

// UserStats represents aggregated discussion metrics for a user profile.
type UserStats struct {
	TopicsCount    int `json:"topics_count"`
	RepliesCount   int `json:"replies_count"`
	SolutionsCount int `json:"solutions_count"`
	TotalVotes     int `json:"total_votes"`
}

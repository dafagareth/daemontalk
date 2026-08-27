package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"daemontalk/internal/i18n"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
)

type terminalPostItem struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Lang        string   `json:"lang"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Cover       string   `json:"cover,omitempty"`
	MinRead     int      `json:"min_read"`
	BodySnippet string   `json:"body_snippet"`
}

type terminalProjectItem struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Role        string   `json:"role"`
	Year        string   `json:"year"`
	Tech        []string `json:"tech"`
	URL         string   `json:"url"`
}

type terminalDataResponse struct {
	Host     string                `json:"host"`
	User     string                `json:"user"`
	Version  string                `json:"version"`
	Posts    []terminalPostItem    `json:"posts"`
	Projects []terminalProjectItem `json:"projects"`
	Tags     map[string]int        `json:"tags"`
	Bio      string                `json:"bio"`
	Socials  map[string]string     `json:"socials"`
}

// Terminal renders the interactive retro in-browser terminal page.
func (h *Handler) Terminal(w http.ResponseWriter, r *http.Request) {
	lang := "en"
	if strings.HasPrefix(r.URL.Path, "/id") {
		lang = "id"
	}
	ui := i18n.Get(lang)

	desc := "Interactive retro UNIX-style web console. Browse articles, run Go code, explore system tools and projects directly in your browser."
	if lang == "id" {
		desc = "Konsol retro UNIX di browser. Telusuri artikel, eksekusi kode Go, jelajahi proyek dan tools langsung di web."
	}

	meta := templates.PageMeta{
		Description: desc,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.Terminal(ui, lang, "terminal", r.URL.Path, meta).Render(r.Context(), w)
}

// TerminalData serves structured JSON for rapid local command evaluation in the terminal.
func (h *Handler) TerminalData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")

	var postItems []terminalPostItem
	tagCounts := make(map[string]int)

	for _, p := range h.AllPosts() {
		if p.Draft {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) {
			continue
		}
		// Clean snippet
		plainBody := stripHTML(string(p.Body))
		if len(plainBody) > 350 {
			plainBody = plainBody[:350] + "..."
		}

		postItems = append(postItems, terminalPostItem{
			Slug:        p.Slug,
			Title:       p.Title,
			Date:        p.Date.Format("2006-01-02"),
			Lang:        p.Lang,
			Tags:        p.Tags,
			Description: p.Description,
			Cover:       p.Cover,
			MinRead:     p.ReadTime,
			BodySnippet: plainBody,
		})

		for _, t := range p.Tags {
			tagCounts[strings.ToLower(t)]++
		}
	}

	var projectItems []terminalProjectItem
	for _, pr := range project.All {
		projectItems = append(projectItems, terminalProjectItem{
			Slug:        pr.Slug,
			Title:       pr.Name,
			Description: pr.Description,
			Tech:        pr.TechStack,
			URL:         pr.RepoURL,
		})
	}

	resp := terminalDataResponse{
		Host:     "daemontalk.local",
		User:     "visitor",
		Version:  "2.8.0-lts",
		Posts:    postItems,
		Projects: projectItems,
		Tags:     tagCounts,
		Bio:      "daemontalk - Open engineering notebook & study log. Exploring Go backend, Python, Linux systems, and backend architecture.",
		Socials: map[string]string{
			"github":    "https://github.com/dafagareth/daemontalk",
			"x":         "https://x.com/daemontalk",
			"instagram": "https://instagram.com/daemontalk",
			"threads":   "https://threads.com/daemontalk",
			"email":     "realdaemontalk@gmail.com",
			"rss":       "/rss.xml",
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func stripHTML(input string) string {
	var sb strings.Builder
	inTag := false
	for _, r := range input {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			sb.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(sb.String()), " ")
}

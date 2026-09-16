package router

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"daemontalk/internal/handler"
)

func redirect301(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handler.SecurityHeaders)
	r.Use(h.Analytics)
	r.Use(h.AuthMiddleware)
	r.NotFound(h.NotFound)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})
	r.With(handler.StaticCacheControl).HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/logo/favicon.ico")
	})
	r.HandleFunc("/google{code}.html", func(w http.ResponseWriter, r *http.Request) {
		file := filepath.Base(r.URL.Path)
		if strings.HasPrefix(file, "google") && strings.HasSuffix(file, ".html") {
			http.ServeFile(w, r, file)
			return
		}
		http.NotFound(w, r)
	})
	r.Get("/og.png", h.SiteOGImage)
	r.Head("/og.png", h.SiteOGImage)

	commentLimit := handler.NewRateLimiter(5, time.Minute)
	contactLimit := handler.NewRateLimiter(3, time.Hour)
	reactionLimit := handler.NewRateLimiter(10, time.Minute)
	adminLimit := handler.NewRateLimiter(30, time.Minute)

	r.With(handler.StaticCacheControl).Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.Get("/", h.BlogIndex)
	r.Get("/projects", redirect301("/colophon"))
	r.Get("/blog", func(w http.ResponseWriter, r *http.Request) {
		if tag := strings.TrimSpace(r.URL.Query().Get("tag")); tag != "" {
			http.Redirect(w, r, "/blog/tag/"+url.PathEscape(strings.ToLower(tag)), http.StatusMovedPermanently)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	r.Get("/blog/{slug}", h.BlogPost)
	r.Get("/blog/{slug}/og.png", h.OGImage)
	r.Head("/blog/{slug}/og.png", h.OGImage)
	r.Get("/blog/{slug}/comments/stream", h.StreamComments)
	r.Get("/blog/{slug}/comments", h.CommentsPartial)
	r.With(commentLimit).Post("/blog/{slug}/comments", h.PostComment)
	r.Post("/blog/{slug}/comments/{id}/delete", h.DeleteComment)
	r.Get("/blog/{slug}/comments/{id}/edit", h.EditCommentForm)
	r.Post("/blog/{slug}/comments/{id}/update", h.UpdateComment)
	r.Post("/blog/{slug}/comments/{id}/report", h.ReportComment)
	r.Get("/graph", redirect301("/"))
	r.Get("/about", h.About)
	r.Get("/colophon", h.Colophon)
	r.Get("/blog/tag", h.RedirectTag)
	r.Get("/blog/tag/{tag}", h.TagIndex)
	r.Get("/blog/posts", h.BlogPostsPartial)
	r.Get("/blog/tag-posts", h.TagPostsPartial)
	r.With(reactionLimit).Post("/blog/{slug}/reactions/{emoji}", h.PostReaction)
	r.Get("/terminal", redirect301("/"))
	r.Get("/rss.xml", h.RSS)
	r.Get("/feed.xml", h.RSS)
	r.Get("/feed.json", h.JSONFeed)
	r.Get("/sitemap.xml", h.Sitemap)
	r.Get("/sitemap-index.xml", h.SitemapIndex)
	r.Get("/sitemap-en.xml", h.SitemapEN)
	r.Get("/sitemap-id.xml", h.SitemapID)
	r.Get("/robots.txt", h.Robots)
	r.Get("/manifest.json", h.Manifest)
	r.With(contactLimit).Post("/contact", h.Contact)
	r.Get("/contact", h.ContactPage)
	r.Get("/uses", redirect301("/colophon"))
	r.Get("/now", redirect301("/colophon"))
	r.Get("/saved", h.Saved)
	r.Get("/stats", h.Stats)
	r.Get("/search", h.Search)

	r.Get("/auth/github", h.AuthGitHub)
	r.Get("/auth/github/callback", h.AuthGitHubCallback)
	r.Post("/auth/logout", h.AuthLogout)
	r.Get("/auth/logout", h.AuthLogout)
	r.Get("/auth/export", h.AuthExport)
	r.Post("/auth/delete-account", h.AuthDeleteAccount)
	r.Get("/auth/me", h.AuthMe)
	r.Get("/auth/badge", h.AuthBadge)
	r.Get("/settings", h.AuthSettings)
	r.Get("/profile", h.AuthMyProfile)
	r.Get("/u/{username}", h.AuthUserProfile)

	r.Get("/socket", h.Discussions)
	r.Get("/socket/new", h.DiscussionsNew)
	r.With(commentLimit).Post("/socket/new", h.DiscussionsCreate)
	r.Get("/socket/{slug}", h.DiscussionsDetail)
	r.With(commentLimit).Post("/socket/{id}/reply", h.DiscussionsReply)
	r.Post("/socket/{id}/solve", h.DiscussionsSolve)
	r.Post("/socket/{id}/delete", h.DiscussionsDeleteTopic)
	r.Post("/socket/reply/{id}/delete", h.DiscussionsDeleteReply)
	r.With(reactionLimit).Post("/socket/vote", h.DiscussionsVote)

	r.Get("/discussions", redirect301("/socket"))
	r.Get("/discussions/*", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/socket"+strings.TrimPrefix(r.URL.Path, "/discussions"), http.StatusMovedPermanently)
	})
	r.Get("/guestbook", redirect301("/socket"))

	r.Get("/resume", redirect301("/about"))
	r.Get("/changelog", redirect301("/colophon"))
	r.Get("/contribute", h.Contribute)
	r.Get("/daemontalk-template.md", h.DownloadTemplate)
	r.Get("/download/*", redirect301("/daemontalk-template.md"))
	r.Get("/template.md", redirect301("/daemontalk-template.md"))
	r.Get("/privacy", h.Privacy)
	r.Get("/terms", h.Terms)
	r.Get("/license", h.License)
	r.Get("/accessibility", h.Accessibility)
	r.Get("/projects/{slug}", redirect301("/colophon"))
	r.Get("/api/webhook/github", h.GitHubWebhook)
	r.Post("/api/webhook/github", h.GitHubWebhook)

	r.Route("/admin", func(r chi.Router) {
		r.Use(adminLimit)
		r.Get("/", h.Admin)
		r.Get("/dashboard", h.Admin)
		r.Get("/content", h.Admin)
		r.Get("/analytics", h.Admin)
		r.Get("/comments", h.Admin)
		r.Get("/digest", h.Admin)
		r.Post("/comments/{id}/delete", h.AdminDeleteComment)
		r.Get("/posts/new", redirect301("/admin/content"))
		r.Get("/post/new", redirect301("/admin/content"))
		r.Get("/new", redirect301("/admin/content"))
		r.Get("/posts/{id}/edit", redirect301("/admin/content"))
		r.Get("/post/{id}/edit", redirect301("/admin/content"))
		r.Post("/posts/upload-md", h.AdminPostUploadMD)
		r.Post("/upload-image", h.AdminUploadImage)
		r.Get("/posts/file-edit", h.AdminPostFileEdit)
		r.Post("/posts/file-save", h.AdminPostFileSave)
		r.Get("/posts/export", h.AdminPostExportMD)
		r.Post("/posts/file-archive", h.AdminPostFileArchive)
		r.Post("/posts/file-restore", h.AdminPostFileRestore)
		r.Post("/posts/file-delete", h.AdminPostFileDelete)
	})

	r.HandleFunc("/id", redirect301("/"))
	r.HandleFunc("/id/*", func(w http.ResponseWriter, r *http.Request) {
		target := strings.TrimPrefix(r.URL.Path, "/id")
		if target == "" {
			target = "/"
		}
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	return r
}

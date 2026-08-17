package router

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"daemontalk/internal/handler"
)

// redirect301 returns a handler that forwards the old route to a new location
// permanently (SEO: old links stay alive).
func redirect301(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

// New creates and configures the main HTTP router
func New(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handler.SecurityHeaders)
	r.Use(h.Analytics)
	r.NotFound(h.NotFound)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})
	r.With(handler.StaticCacheControl).HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/logo/favicon.ico")
	})
	r.HandleFunc("/google*.html", func(w http.ResponseWriter, r *http.Request) {
		file := filepath.Base(r.URL.Path)
		if strings.HasPrefix(file, "google") && strings.HasSuffix(file, ".html") {
			http.ServeFile(w, r, file)
			return
		}
		http.NotFound(w, r)
	})
	r.Get("/og.png", h.SiteOGImage)

	// Per-IP rate limits
	commentLimit := handler.NewRateLimiter(5, time.Minute)
	contactLimit := handler.NewRateLimiter(3, time.Hour)
	reactionLimit := handler.NewRateLimiter(10, time.Minute)
	runLimit := handler.NewRateLimiter(25, time.Minute)
	adminLimit := handler.NewRateLimiter(30, time.Minute) // Protect against brute-forcing the ADMIN_TOKEN

	r.With(handler.StaticCacheControl).Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.With(handler.StaticCacheControl).Handle("/id/static/*", http.StripPrefix("/id/static/", http.FileServer(http.Dir("web/static"))))

	r.Get("/", h.BlogIndex)
	r.Get("/projects", redirect301("/behind#projects"))
	r.Get("/blog", redirect301("/"))
	r.Get("/blog/{slug}", h.BlogPost)
	r.Get("/blog/{slug}/og.png", h.OGImage)
	r.Get("/blog/{slug}/comments/stream", h.StreamComments)
	r.Get("/blog/{slug}/comments", h.CommentsPartial)
	r.With(commentLimit).Post("/blog/{slug}/comments", h.PostComment)
	r.Post("/blog/{slug}/comments/{id}/delete", h.DeleteComment)
	r.Get("/graph", h.Graph)
	r.Get("/install.sh", h.InstallScript)
	r.Get("/about", h.About)
	r.Get("/behind", h.Behind)
	r.Get("/blog/tag/{tag}", h.TagIndex)
	r.Get("/blog/posts", h.BlogPostsPartial)
	r.With(reactionLimit).Post("/blog/{slug}/reactions/{emoji}", h.PostReaction)
	r.Get("/terminal", h.Terminal)
	r.Get("/api/terminal/data", h.TerminalData)
	r.With(runLimit).Post("/api/run", h.RunCode)
	r.Get("/daily", h.CLIDaily)
	r.Get("/t/{tag}", h.CLITag)
	r.Get("/p/{slug}", h.CLIPost)
	r.Get("/recipes", h.CLIRecipes)
	r.Get("/cheat", h.CLIRecipes)
	r.Get("/ebpf", h.CLIRecipes)
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
	r.Get("/uses", redirect301("/behind"))
	r.Get("/now", redirect301("/behind"))
	r.Get("/saved", h.Saved)
	r.Get("/stats", h.Stats)
	r.Get("/search", h.Search)
	r.Get("/guestbook", h.Guestbook)
	r.With(commentLimit).Post("/guestbook", h.PostGuestbook)
	r.Get("/til", h.TIL)
	r.Get("/resume", h.Resume)
	r.Get("/changelog", h.Changelog)
	r.Get("/contribute", h.Contribute)
	r.Get("/download/template.md", h.DownloadTemplate)
	r.Get("/template.md", h.DownloadTemplate)
	r.Get("/links", h.Links)
	r.Get("/privacy", h.Privacy)
	r.Get("/terms", h.Terms)
	r.Get("/license", h.License)
	r.Get("/projects/{slug}", h.ProjectDetail)

	// Admin routes (Protected by AdminToken and Rate Limited)
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminLimit)
		r.Get("/", h.Admin)
		r.Get("/dashboard", h.Admin)
		r.Get("/content", h.Admin)
		r.Get("/analytics", h.Admin)
		r.Get("/comments", h.Admin)
		r.Get("/digest", h.Admin)
		r.Post("/comments/{id}/delete", h.AdminDeleteComment)
		r.Get("/posts/new", h.AdminPostNew)
		r.Get("/post/new", redirect301("/admin/posts/new"))
		r.Get("/new", redirect301("/admin/posts/new"))
		r.Get("/posts/{id}/edit", h.AdminPostEdit)
		r.Get("/post/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/admin/posts/"+chi.URLParam(r, "id")+"/edit", http.StatusMovedPermanently)
		})
		r.Post("/posts/autosave", h.AdminPostAutosave)
		r.Post("/posts/{id}/publish", h.AdminPostPublish)
		r.Post("/posts/{id}/delete", h.AdminPostDelete)
		r.Post("/posts/upload-md", h.AdminPostUploadMD)
		r.Post("/upload-image", h.AdminUploadImage)
		r.Get("/posts/file-edit", h.AdminPostFileEdit)
		r.Post("/posts/file-save", h.AdminPostFileSave)
		r.Get("/posts/export", h.AdminPostExportMD)
		r.Post("/posts/file-archive", h.AdminPostFileArchive)
		r.Post("/posts/file-restore", h.AdminPostFileRestore)
		r.Post("/posts/file-delete", h.AdminPostFileDelete)
		r.Post("/settings/toggle-radar", h.AdminToggleRadar)
	})

	r.Route("/id", func(r chi.Router) {
		r.Get("/", h.BlogIndex)
		r.Get("/projects", redirect301("/id/behind#projects"))
		r.Get("/blog", redirect301("/id"))
		r.Get("/graph", h.Graph)
		r.Get("/blog/tag/{tag}", h.TagIndex)
		r.Get("/blog/{slug}", h.BlogPost)
		r.Get("/blog/{slug}/comments/stream", h.StreamComments)
		r.Get("/blog/{slug}/comments", h.CommentsPartial)
		r.With(commentLimit).Post("/blog/{slug}/comments", h.PostComment)
		r.Post("/blog/{slug}/comments/{id}/delete", h.DeleteComment)
		r.With(reactionLimit).Post("/blog/{slug}/reactions/{emoji}", h.PostReaction)
		r.Get("/terminal", h.Terminal)
		r.Get("/about", h.About)
		r.Get("/behind", h.Behind)
		r.Get("/uses", redirect301("/id/behind"))
		r.Get("/now", redirect301("/id/behind"))
		r.Get("/saved", h.Saved)
		r.Get("/stats", h.Stats)
		r.Get("/daily", h.CLIDaily)
		r.Get("/search", h.Search)
		r.Get("/guestbook", h.Guestbook)
		r.With(commentLimit).Post("/guestbook", h.PostGuestbook)
		r.Get("/til", h.TIL)
		r.Get("/resume", h.Resume)
		r.Get("/changelog", h.Changelog)
		r.Get("/contribute", h.Contribute)
		r.Get("/download/template.md", h.DownloadTemplate)
		r.Get("/template.md", h.DownloadTemplate)
		r.Get("/links", h.Links)
		r.Get("/privacy", h.Privacy)
		r.Get("/terms", h.Terms)
		r.Get("/license", h.License)
		r.Get("/projects/{slug}", h.ProjectDetail)
	})

	return r
}

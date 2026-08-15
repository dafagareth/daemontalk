package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"daemontalk/internal/comment"
	"daemontalk/internal/handler"
	"daemontalk/internal/highlight"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
	"daemontalk/internal/tuisrv"
	"daemontalk/web/templates"
)

// redirect301 mengembalikan handler yang meneruskan route lama ke lokasi baru
// secara permanen (SEO: link lama tetap hidup setelah pivot blog-first).
func redirect301(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

func main() {
	var logHandler slog.Handler
	if os.Getenv("ENV") == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(logHandler))

	posts, err := post.LoadAllWithDrafts("content/posts")
	if err != nil {
		slog.Error("load posts failed", "error", err)
		os.Exit(1)
	}

	if err := os.WriteFile("web/static/css/chroma.css", []byte(highlight.GenerateCSS()), 0644); err != nil {
		slog.Warn("chroma.css generation failed", "error", err)
	}

	// Cache-bust static CSS: use main.css modtime so a rebuild invalidates the
	// browser cache automatically.
	if fi, err := os.Stat("web/static/css/main.css"); err == nil {
		templates.AssetVersion = fmt.Sprintf("%d", fi.ModTime().Unix())
	}
	for _, sz := range []int{192, 512} {
		name := fmt.Sprintf("web/static/images/icon-%d.png", sz)
		if _, err := os.Stat(name); os.IsNotExist(err) {
			if err := os.WriteFile(name, generateIcon(sz), 0644); err != nil {
				slog.Warn("failed to write app icon", "size", sz, "path", name, "error", err)
			}
		}
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		slog.Error("create data dir failed", "error", err)
		os.Exit(1)
	}
	comments, err := comment.Open("data/comments.db")
	if err != nil {
		slog.Error("open comments db failed", "error", err)
		os.Exit(1)
	}
	defer comments.Close()

	// Post buatan editor web — persisten di volume data/ bersama comments.db.
	pdb, err := postdb.Open("data/posts.db")
	if err != nil {
		slog.Error("open posts db failed", "error", err)
		os.Exit(1)
	}
	defer pdb.Close()

	h := &handler.Handler{
		AllProjects: project.All,
		FilePosts:   posts,
		PostDB:      pdb,
		Comments:    comments,
		AdminToken:  os.Getenv("ADMIN_TOKEN"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		SMTPTo:      os.Getenv("SMTP_TO"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
	}
	h.RefreshPosts()
	if h.AdminToken != "" {
		slog.Info("comment moderation enabled", "login_via", "?admin=TOKEN")
	}
	if h.SMTPHost != "" {
		slog.Info("contact form SMTP enabled", "host", h.SMTPHost)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handler.SecurityHeaders)
	r.Use(h.Analytics)
	r.NotFound(h.NotFound)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/logo/favicon.ico")
	})
	r.Get("/og.png", h.SiteOGImage)

	// Per-IP rate limits
	commentLimit  := handler.NewRateLimiter(5, time.Minute)
	contactLimit  := handler.NewRateLimiter(3, time.Hour)
	reactionLimit := handler.NewRateLimiter(10, time.Minute)
	runLimit      := handler.NewRateLimiter(25, time.Minute)
	adminLimit    := handler.NewRateLimiter(30, time.Minute) // Protect against brute-forcing the ADMIN_TOKEN

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.Handle("/id/static/*", http.StripPrefix("/id/static/", http.FileServer(http.Dir("web/static"))))

	r.Get("/", h.BlogIndex)
	r.Get("/projects", redirect301("/behind#projects"))
	r.Get("/blog", redirect301("/"))
	r.Get("/blog/{slug}", h.BlogPost)
	r.Get("/blog/{slug}/og.png", h.OGImage)
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
	})

	r.Route("/id", func(r chi.Router) {
		r.Get("/", h.BlogIndex)
		r.Get("/projects", redirect301("/id/behind#projects"))
		r.Get("/blog", redirect301("/id"))
		r.Get("/graph", h.Graph)
		r.Get("/blog/tag/{tag}", h.TagIndex)
		r.Get("/blog/{slug}", h.BlogPost)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Listen for OS termination signals so the server can drain in-flight
	// requests before exiting (important under Docker / orchestrators).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start Wish SSH Server for direct TUI access (ssh daemontalk.com -p 2222)
	sshPort := os.Getenv("SSH_PORT")
	if sshPort == "" {
		sshPort = "2222"
	}
	sshHostKey := filepath.Join("data", ".ssh_host_key")
	sshSrv, err := tuisrv.Start(":"+sshPort, sshHostKey)
	if err != nil {
		slog.Warn("failed to initialize SSH TUI server", "error", err)
	} else {
		go func() {
			slog.Info("SSH TUI server starting", "port", sshPort)
			if err := sshSrv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				slog.Error("SSH TUI server failed", "error", err)
			}
		}()
	}

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down servers...")

	if sshSrv != nil {
		sshShutdownCtx, sshCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer sshCancel()
		_ = sshSrv.Shutdown(sshShutdownCtx)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
	slog.Info("server stopped")
}

// generateIcon creates a minimal solid-color square PNG icon for the PWA manifest.
func generateIcon(size int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	blue := color.RGBA{0x1a, 0x73, 0xe8, 0xff}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

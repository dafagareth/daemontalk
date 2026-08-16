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
	"daemontalk/internal/comment"
	"daemontalk/internal/handler"
	"daemontalk/internal/highlight"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
	"daemontalk/internal/router"
	"daemontalk/internal/tuisrv"
	"daemontalk/web/templates"
)


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

	r := router.New(h)
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

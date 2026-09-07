package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type githubPushPayload struct {
	Ref string `json:"ref"`
}

func (h *Handler) GitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"message": "DaemonTalk GitHub webhook endpoint is active",
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if sig == "" {
			http.Error(w, "Missing X-Hub-Signature-256 header", http.StatusUnauthorized)
			return
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
			http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
			return
		}
	}

	event := r.Header.Get("X-GitHub-Event")
	if event == "ping" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "message": "pong"})
		return
	}

	var payload githubPushPayload
	if err := json.Unmarshal(body, &payload); err == nil {
		isMain := strings.HasSuffix(payload.Ref, "/main") || strings.HasSuffix(payload.Ref, "/master")
		isTag := strings.HasPrefix(payload.Ref, "refs/tags/")
		if payload.Ref != "" && !isMain && !isTag {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status":  "ignored",
				"message": "Push is not to main branch or release tag (" + payload.Ref + ")",
			})
			return
		}
	}

	h.ReloadFilePosts()
	h.RefreshPosts()

	slog.Info("github webhook: successfully reloaded dispatches and refreshed post cache")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "dispatches successfully reloaded",
		"total":   len(h.AllPosts()),
	})
}

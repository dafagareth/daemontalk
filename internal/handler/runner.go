package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"daemontalk/internal/runner"
)

type runRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// RunCode handles code execution requests for in-article snippets and terminal.
func (h *Handler) RunCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Limit request payload to 64KB
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(runner.RunResponse{
			Error:   "invalid request payload",
			Success: false,
		})
		return
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		_ = json.NewEncoder(w).Encode(runner.RunResponse{
			Error:   "empty code payload",
			Success: false,
		})
		return
	}

	lang := strings.ToLower(strings.TrimSpace(req.Language))
	start := time.Now()

	switch lang {
	case "go", "golang":
		resp := runner.ExecuteGoCode(r.Context(), code)
		resp.DurationMs = time.Since(start).Milliseconds()
		_ = json.NewEncoder(w).Encode(resp)

	case "python", "py", "python3":
		resp := runner.ExecutePythonCode(r.Context(), code)
		resp.DurationMs = time.Since(start).Milliseconds()
		_ = json.NewEncoder(w).Encode(resp)

	case "sh", "bash", "shell":
		resp := runner.ExecuteShellSim(r.Context(), code)
		resp.DurationMs = time.Since(start).Milliseconds()
		_ = json.NewEncoder(w).Encode(resp)

	case "js", "javascript", "node":
		resp := runner.ExecuteNodeCode(r.Context(), code)
		resp.DurationMs = time.Since(start).Milliseconds()
		_ = json.NewEncoder(w).Encode(resp)

	default:
		_ = json.NewEncoder(w).Encode(runner.RunResponse{
			Error:      fmt.Sprintf("unsupported language: %s (supported: go, js, python, bash)", lang),
			DurationMs: time.Since(start).Milliseconds(),
			Success:    false,
		})
	}
}

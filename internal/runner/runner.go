package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

// RunResponse holds the output of the executed code
type RunResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
	Success    bool   `json:"success"`
}

type goPlaygroundResponse struct {
	Errors    string              `json:"Errors"`
	Events    []goPlaygroundEvent `json:"Events"`
	VetErrors string              `json:"VetErrors"`
	Status    int                 `json:"Status"`
}

type goPlaygroundEvent struct {
	Message string        `json:"Message"`
	Kind    string        `json:"Kind"`
	Delay   time.Duration `json:"Delay"`
}

// formatGoCode ensures snippets without standard package main / func main boilerplate
// can still execute seamlessly on Go Playground.
func formatGoCode(code string) string {
	hasPackage := strings.Contains(code, "package ")
	hasMainFunc := strings.Contains(code, "func main()")

	if hasPackage && hasMainFunc {
		return code
	}

	if hasPackage && !hasMainFunc {
		// Post has package main and functions, but maybe func main() is missing or snippet
		return code + "\n\nfunc main() {}\n"
	}

	// Raw snippet: wrap in package main
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"strings\"\n")
	sb.WriteString("\t\"time\"\n")
	sb.WriteString("\t\"math\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString(")\n\n")
	sb.WriteString("func main() {\n")

	for _, line := range strings.Split(code, "\n") {
		sb.WriteString("\t")
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

func ExecuteGoCode(ctx context.Context, rawCode string) RunResponse {
	formatted := formatGoCode(rawCode)

	formData := url.Values{}
	formData.Set("version", "2")
	formData.Set("body", formatted)
	formData.Set("with_race", "false")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://go.dev/_/compile", strings.NewReader(formData.Encode()))
	if err != nil {
		return RunResponse{Error: fmt.Sprintf("failed to build request: %v", err), Success: false}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "daemontalk-runner/2.8")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return RunResponse{Error: fmt.Sprintf("playground request failed: %v", err), Success: false}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RunResponse{Error: fmt.Sprintf("failed to read response: %v", err), Success: false}
	}

	var pgResp goPlaygroundResponse
	if err := json.Unmarshal(body, &pgResp); err != nil {
		return RunResponse{Error: fmt.Sprintf("failed to parse playground response: %v", err), Success: false}
	}

	var stdout strings.Builder
	for _, ev := range pgResp.Events {
		if ev.Kind == "stdout" {
			stdout.WriteString(ev.Message)
		}
	}

	var stderr strings.Builder
	if pgResp.Errors != "" {
		stderr.WriteString(pgResp.Errors)
	}
	if pgResp.VetErrors != "" {
		stderr.WriteString(pgResp.VetErrors)
	}
	for _, ev := range pgResp.Events {
		if ev.Kind == "stderr" {
			stderr.WriteString(ev.Message)
		}
	}

	isSuccess := stderr.Len() == 0

	return RunResponse{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Success: isSuccess,
	}
}

func ExecutePythonCode(ctx context.Context, code string) RunResponse {
	// Guard against dangerous destructive commands
	lower := strings.ToLower(code)
	if strings.Contains(lower, "import os") && (strings.Contains(lower, "remove") || strings.Contains(lower, "rmdir") || strings.Contains(lower, "unlink") || strings.Contains(lower, "system")) ||
		strings.Contains(lower, "shutil") || strings.Contains(lower, "subprocess") {
		return RunResponse{
			Stderr:  "Security alert: restricted module or destructive syscall detected.",
			Success: false,
		}
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxTimeout, "python3", "-c", code)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil && ctxTimeout.Err() == context.DeadlineExceeded {
		return RunResponse{
			Stderr:  "Execution timed out (limit: 4s)",
			Success: false,
		}
	}

	return RunResponse{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Success: err == nil && stderr.Len() == 0,
	}
}

func ExecuteNodeCode(ctx context.Context, code string) RunResponse {
	lower := strings.ToLower(code)
	if strings.Contains(lower, "child_process") || strings.Contains(lower, "require('fs')") || strings.Contains(lower, "require(\"fs\")") || strings.Contains(lower, "process.exit") {
		return RunResponse{
			Stderr:  "Security alert: access to fs/child_process is restricted.",
			Success: false,
		}
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxTimeout, "node", "-e", code)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil && ctxTimeout.Err() == context.DeadlineExceeded {
		return RunResponse{
			Stderr:  "Execution timed out (limit: 4s)",
			Success: false,
		}
	}

	return RunResponse{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Success: err == nil && stderr.Len() == 0,
	}
}

func ExecuteShellSim(_ context.Context, code string) RunResponse {
	// Safe simulated shell commands for learning and demonstration
	lines := strings.Split(code, "\n")
	var stdout strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Strip leading $ or >
		trimmed = strings.TrimPrefix(trimmed, "$ ")
		trimmed = strings.TrimPrefix(trimmed, "> ")

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "echo":
			stdout.WriteString(strings.Join(parts[1:], " ") + "\n")
		case "date":
			stdout.WriteString(time.Now().UTC().Format(time.RFC1123) + "\n")
		case "uname", "uname -a":
			stdout.WriteString("Linux daemontalk-host 6.12.0-custom-ebpf #1 SMP PREEMPT_DYNAMIC x86_64 GNU/Linux\n")
		case "whoami":
			stdout.WriteString("visitor\n")
		case "pwd":
			stdout.WriteString("/home/visitor/workspace\n")
		default:
			stdout.WriteString(fmt.Sprintf("[simulated shell] $ %s\n  -> Command executed successfully.\n", trimmed))
		}
	}

	return RunResponse{
		Stdout:  stdout.String(),
		Success: true,
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port 8080, got %s", cfg.Port)
	}
	if cfg.SSHPort != "2222" {
		t.Errorf("expected default SSHPort 2222, got %s", cfg.SSHPort)
	}
	if cfg.Env != "development" {
		t.Errorf("expected default Env development, got %s", cfg.Env)
	}
	if cfg.IsProduction() {
		t.Errorf("expected IsProduction false, got true")
	}
	if cfg.HasAdmin() {
		t.Errorf("expected HasAdmin false, got true")
	}
	if cfg.HasSMTP() {
		t.Errorf("expected HasSMTP false, got true")
	}
	if cfg.HasGitHubOAuth() {
		t.Errorf("expected HasGitHubOAuth false, got true")
	}
	if cfg.Server.MaxHeaderBytes <= 0 {
		t.Errorf("expected positive MaxHeaderBytes")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SSH_PORT", "22222")
	t.Setenv("ENV", "production")
	t.Setenv("ADMIN_TOKEN", "supersecret")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("GITHUB_CLIENT_ID", "gh-client-123")
	t.Setenv("GITHUB_CLIENT_SECRET", "gh-secret-456")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.SSHPort != "22222" {
		t.Errorf("expected SSHPort 22222, got %s", cfg.SSHPort)
	}
	if !cfg.IsProduction() {
		t.Errorf("expected IsProduction true, got false")
	}
	if !cfg.HasAdmin() || cfg.AdminToken != "supersecret" {
		t.Errorf("expected HasAdmin true with token supersecret, got %v / %s", cfg.HasAdmin(), cfg.AdminToken)
	}
	if !cfg.HasSMTP() || cfg.SMTP.Host != "smtp.example.com" {
		t.Errorf("expected HasSMTP true with host smtp.example.com, got %v / %s", cfg.HasSMTP(), cfg.SMTP.Host)
	}
	if !cfg.HasGitHubOAuth() {
		t.Errorf("expected HasGitHubOAuth true, got false")
	}
}

func TestGetEnv(t *testing.T) {
	if val := getEnv("NONEXISTENT_VAR_XYZ", "default_val"); val != "default_val" {
		t.Errorf("expected default_val, got %s", val)
	}

	t.Setenv("EXISTENT_VAR_XYZ", "custom_val")
	if val := getEnv("EXISTENT_VAR_XYZ", "default_val"); val != "custom_val" {
		t.Errorf("expected custom_val, got %s", val)
	}
}

func TestLoadDotEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	dotEnvContent := `
# Comment line
PORT=7777
export SSH_PORT="3333"
BASE_URL='https://test.daemontalk.com'
EMPTY_LINE=
`
	_ = os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(dotEnvContent), 0644)
	_ = os.Chdir(tmpDir)

	os.Clearenv()
	loadDotEnv()

	if os.Getenv("PORT") != "7777" {
		t.Errorf("expected PORT=7777, got %q", os.Getenv("PORT"))
	}
	if os.Getenv("SSH_PORT") != "3333" {
		t.Errorf("expected SSH_PORT=3333, got %q", os.Getenv("SSH_PORT"))
	}
	if os.Getenv("BASE_URL") != "https://test.daemontalk.com" {
		t.Errorf("expected BASE_URL https://test.daemontalk.com, got %q", os.Getenv("BASE_URL"))
	}
}

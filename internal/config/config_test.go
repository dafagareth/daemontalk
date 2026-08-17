package config

import (
	"os"
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
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SSH_PORT", "22222")
	t.Setenv("ENV", "production")
	t.Setenv("ADMIN_TOKEN", "supersecret")
	t.Setenv("SMTP_HOST", "smtp.example.com")

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
}

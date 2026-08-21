package config

import (
	"os"
	"strings"
	"time"
)

// Config holds all runtime application configuration loaded from the environment.
type Config struct {
	Port        string
	SSHPort     string
	Env         string
	BaseURL     string
	AdminToken  string
	ContentDir  string
	DataDir     string
	Server      ServerConfig
	SMTP        SMTPConfig
	GitHubToken string
}

// ServerConfig holds HTTP server network timeout and memory limits for resilience.
type ServerConfig struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

// SMTPConfig encapsulates email transmission settings.
type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	To   string
}

// IsProduction returns true if running under production environment mode.
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Env) == "production"
}

// HasSMTP returns true if SMTP host is configured.
func (c *Config) HasSMTP() bool {
	return c.SMTP.Host != ""
}

// HasAdmin returns true if an admin token is set.
func (c *Config) HasAdmin() bool {
	return c.AdminToken != ""
}

// Load reads and validates configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		SSHPort:     getEnv("SSH_PORT", "2222"),
		Env:         getEnv("ENV", "development"),
		BaseURL:     getEnv("BASE_URL", "https://www.daemontalk.com"),
		AdminToken:  getEnv("ADMIN_TOKEN", ""),
		ContentDir:  getEnv("CONTENT_DIR", "content"),
		DataDir:     getEnv("DATA_DIR", "data"),
		Server: ServerConfig{
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MB
		},
		SMTP: SMTPConfig{
			Host: getEnv("SMTP_HOST", ""),
			Port: getEnv("SMTP_PORT", "587"),
			User: getEnv("SMTP_USER", ""),
			Pass: getEnv("SMTP_PASS", ""),
			To:   getEnv("SMTP_TO", ""),
		},
		GitHubToken: getEnv("GITHUB_TOKEN", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

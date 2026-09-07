package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port               string
	SSHPort            string
	Env                string
	BaseURL            string
	AdminToken         string
	ContentDir         string
	DataDir            string
	Server             ServerConfig
	SMTP               SMTPConfig
	GitHubToken        string
	GitHubClientID     string
	GitHubClientSecret string
	SessionSecret      string
}

type ServerConfig struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	To   string
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Env) == "production"
}

func (c *Config) HasSMTP() bool {
	return c.SMTP.Host != ""
}

func (c *Config) HasAdmin() bool {
	return c.AdminToken != ""
}

func (c *Config) HasGitHubOAuth() bool {
	return c.GitHubClientID != "" && c.GitHubClientSecret != ""
}

func Load() *Config {
	loadDotEnv()
	return &Config{
		Port:               getEnv("PORT", "8080"),
		SSHPort:            getEnv("SSH_PORT", "2222"),
		Env:                getEnv("ENV", "development"),
		BaseURL:            getEnv("BASE_URL", "https://www.daemontalk.com"),
		AdminToken:         getEnv("ADMIN_TOKEN", ""),
		ContentDir:         getEnv("CONTENT_DIR", "content"),
		DataDir:            getEnv("DATA_DIR", "data"),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		SessionSecret:      getEnv("SESSION_SECRET", "daemontalk-default-insecure-secret-key-32b"),
		Server: ServerConfig{
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20,
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

func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimSpace(line)

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if val, exists := os.LookupEnv(k); !exists || val == "" || val == "change-me-in-dotenv" {
				_ = os.Setenv(k, v)
			}
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

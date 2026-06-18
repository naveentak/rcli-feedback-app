package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port           string
	GitHubToken    string
	GitHubOwner    string
	GitHubRepo     string
	APIKeys        map[string]string // app name -> API key
	AllowedOrigins []string
	DevMode        bool // local dev: skip API key on all write endpoints
	PublicSubmit   bool // production: allow unauthenticated POST /feedback (web form)
	HMACSecret     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:        envOr("PORT", "8080"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		GitHubOwner: envOr("GITHUB_OWNER", ""),
		GitHubRepo:  envOr("GITHUB_REPO", "rcli-feedback-app"),
		APIKeys:     parseAPIKeys(os.Getenv("API_KEYS")),
		DevMode:      os.Getenv("DEV_MODE") == "true",
		PublicSubmit: os.Getenv("PUBLIC_SUBMIT") == "true",
		HMACSecret:   os.Getenv("FEEDBACK_HMAC_SECRET"),
	}

	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		cfg.AllowedOrigins = strings.Split(origins, ",")
	}

	if cfg.GitHubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required")
	}
	if cfg.GitHubOwner == "" {
		return nil, fmt.Errorf("GITHUB_OWNER is required")
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseAPIKeys parses "rclip:key1,boka:key2" into a map.
func parseAPIKeys(raw string) map[string]string {
	keys := make(map[string]string)
	if raw == "" {
		return keys
	}
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			keys[strings.ToLower(parts[0])] = parts[1]
		}
	}
	return keys
}
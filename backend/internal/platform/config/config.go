// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AdminEmails []string
	BackendDir  string // working dir for the admin "run backend tests" action
	FrontendDir string // working dir for the admin "run frontend tests" action
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://rimu:rimu@localhost:5432/rimu?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		AdminEmails: splitEmails(getEnv("ADMIN_EMAILS", "")),
		BackendDir:  getEnv("BACKEND_DIR", "."),
		FrontendDir: getEnv("FRONTEND_DIR", "../frontend"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitEmails(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	emails := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			emails = append(emails, p)
		}
	}
	return emails
}

func (c Config) IsAdminEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, e := range c.AdminEmails {
		if e == email {
			return true
		}
	}
	return false
}

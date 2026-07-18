// Package testutil boots the real wired server against a dedicated test
// database so every bounded context's e2e tests can drive it purely over
// HTTP, with no mocks.
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/platform/config"
	"rimu/backend/internal/platform/httpserver"
)

type Server struct {
	*httptest.Server
	t *testing.T
}

func NewServer(t *testing.T) *Server {
	t.Helper()

	cfg := config.Config{
		DatabaseURL: testDatabaseURL(),
		JWTSecret:   "test-secret",
		AdminEmails: []string{"admin@rimu.test"},
		BackendDir:  ".",
		FrontendDir: ".",
	}

	router, pool, err := httpserver.Build(context.Background(), cfg)
	if err != nil {
		t.Fatalf("build test server: %v", err)
	}

	truncateAll(t, pool)
	t.Cleanup(pool.Close)

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return &Server{Server: srv, t: t}
}

func testDatabaseURL() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://rimu:rimu@localhost:5432/rimu_test?sslmode=disable"
}

// Do issues an HTTP request against the test server and decodes the JSON
// response body into `out` (pass nil to ignore the body).
func (s *Server) Do(method, path, token string, body, out interface{}) int {
	s.t.Helper()

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			s.t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, s.URL+path, reader)
	if err != nil {
		s.t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil && err != io.EOF {
			s.t.Fatalf("decode response body: %v", err)
		}
	}
	return resp.StatusCode
}

// RegisterAndLogin registers a fresh user and returns their auth token.
func (s *Server) RegisterAndLogin(email, password string) string {
	s.t.Helper()

	status := s.Do(http.MethodPost, "/api/auth/register", "", map[string]string{
		"email": email, "password": password,
	}, nil)
	if status != http.StatusCreated {
		s.t.Fatalf("register %s: unexpected status %d", email, status)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	status = s.Do(http.MethodPost, "/api/auth/login", "", map[string]string{
		"email": email, "password": password,
	}, &loginResp)
	if status != http.StatusOK {
		s.t.Fatalf("login %s: unexpected status %d", email, status)
	}
	return loginResp.Token
}

// Upgrade flips the given account to the pro plan via the mock endpoint.
func (s *Server) Upgrade(token string) {
	s.t.Helper()
	if status := s.Do(http.MethodPost, "/api/me/upgrade", token, nil, nil); status != http.StatusOK {
		s.t.Fatalf("upgrade: unexpected status %d", status)
	}
}

func truncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE users, groups, group_members, habits, habit_logs,
			workout_sessions, workout_sets, transactions, notes, roadmap_items
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestUser_RegisterLoginMeUpgrade(t *testing.T) {
	s := testutil.NewServer(t)

	status := s.Do(http.MethodPost, "/api/auth/register", "", map[string]string{
		"email": "alice@example.com", "password": "password123",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", status)
	}

	// Duplicate email is rejected.
	status = s.Do(http.MethodPost, "/api/auth/register", "", map[string]string{
		"email": "alice@example.com", "password": "password123",
	}, nil)
	if status != http.StatusConflict {
		t.Fatalf("duplicate register: expected 409, got %d", status)
	}

	token := s.RegisterAndLogin("bob@example.com", "password123")
	if token == "" {
		t.Fatal("expected a non-empty token")
	}

	var me struct {
		Email string `json:"email"`
		Plan  string `json:"plan"`
	}
	status = s.Do(http.MethodGet, "/api/me", token, nil, &me)
	if status != http.StatusOK || me.Email != "bob@example.com" || me.Plan != "free" {
		t.Fatalf("me: unexpected response status=%d body=%+v", status, me)
	}

	var upgraded struct {
		Plan string `json:"plan"`
	}
	status = s.Do(http.MethodPost, "/api/me/upgrade", token, nil, &upgraded)
	if status != http.StatusOK || upgraded.Plan != "pro" {
		t.Fatalf("upgrade: expected plan=pro, got status=%d body=%+v", status, upgraded)
	}
}

func TestUser_MockUpgradeRefusedWhenBillingConfigured(t *testing.T) {
	s := testutil.NewServerWithStripeConfigured(t)
	token := s.RegisterAndLogin("carol@example.com", "password123")

	status := s.Do(http.MethodPost, "/api/me/upgrade", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("mock upgrade with real billing configured: expected 403, got %d", status)
	}

	var me struct {
		Plan string `json:"plan"`
	}
	status = s.Do(http.MethodGet, "/api/me", token, nil, &me)
	if status != http.StatusOK || me.Plan != "free" {
		t.Fatalf("plan should stay free after refused mock upgrade: status=%d body=%+v", status, me)
	}
}

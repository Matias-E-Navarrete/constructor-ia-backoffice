package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestAdmin_AccessControlFlagsUsersRoadmap(t *testing.T) {
	s := testutil.NewServer(t)

	// testutil configures ADMIN_EMAILS=admin@rimu.test, so this registration
	// is auto-promoted to the admin role.
	adminToken := s.RegisterAndLogin("admin@rimu.test", "password123")
	userToken := s.RegisterAndLogin("regular@example.com", "password123")

	status := s.Do(http.MethodGet, "/api/admin/users", userToken, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("non-admin listing users: expected 403, got %d", status)
	}

	// --- feature flags + kill switch ------------------------------------------------
	var flags []struct {
		Key     string `json:"key"`
		Enabled bool   `json:"enabled"`
	}
	status = s.Do(http.MethodGet, "/api/admin/feature-flags", adminToken, nil, &flags)
	if status != http.StatusOK || len(flags) == 0 {
		t.Fatalf("list feature flags: status=%d len=%d", status, len(flags))
	}

	status = s.Do(http.MethodPatch, "/api/admin/feature-flags/habits", adminToken, map[string]bool{"enabled": false}, nil)
	if status != http.StatusOK {
		t.Fatalf("disable habits flag: expected 200, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/habits", userToken, nil, nil)
	if status != http.StatusServiceUnavailable {
		t.Fatalf("habits while disabled: expected 503, got %d", status)
	}

	status = s.Do(http.MethodPatch, "/api/admin/feature-flags/habits", adminToken, map[string]bool{"enabled": true}, nil)
	if status != http.StatusOK {
		t.Fatalf("re-enable habits flag: expected 200, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/habits", userToken, nil, nil)
	if status != http.StatusOK {
		t.Fatalf("habits after re-enable: expected 200, got %d", status)
	}

	// --- user management ---------------------------------------------------------
	var users []struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	status = s.Do(http.MethodGet, "/api/admin/users", adminToken, nil, &users)
	if status != http.StatusOK || len(users) != 2 {
		t.Fatalf("list users: status=%d len=%d", status, len(users))
	}

	var regularUserID string
	for _, u := range users {
		if u.Email == "regular@example.com" {
			regularUserID = u.ID
		}
	}
	if regularUserID == "" {
		t.Fatal("could not find regular user in admin user list")
	}

	var updated struct {
		Plan string `json:"plan"`
	}
	status = s.Do(http.MethodPatch, "/api/admin/users/"+regularUserID+"/plan", adminToken, map[string]string{"plan": "pro"}, &updated)
	if status != http.StatusOK || updated.Plan != "pro" {
		t.Fatalf("set user plan: status=%d body=%+v", status, updated)
	}

	// --- roadmap -------------------------------------------------------------------
	var item struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	status = s.Do(http.MethodPost, "/api/admin/roadmap", adminToken, map[string]string{
		"title": "Agregar notificaciones", "kind": "feature",
	}, &item)
	if status != http.StatusCreated || item.Status != "planned" {
		t.Fatalf("create roadmap item: status=%d body=%+v", status, item)
	}

	status = s.Do(http.MethodPatch, "/api/admin/roadmap/"+item.ID, adminToken, map[string]string{"status": "in_progress"}, &item)
	if status != http.StatusOK || item.Status != "in_progress" {
		t.Fatalf("update roadmap status: status=%d body=%+v", status, item)
	}
}

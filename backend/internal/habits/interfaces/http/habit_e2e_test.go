package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestHabits_CreateCheckInStatsGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var habit struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	status := s.Do(http.MethodPost, "/api/habits", token, map[string]string{"name": "Meditar"}, &habit)
	if status != http.StatusCreated || habit.Name != "Meditar" {
		t.Fatalf("create habit: status=%d body=%+v", status, habit)
	}

	status = s.Do(http.MethodPost, "/api/habits/"+habit.ID+"/log", token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("check-in: expected 204, got %d", status)
	}

	var list []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/habits", token, nil, &list)
	if status != http.StatusOK || len(list) != 1 {
		t.Fatalf("list: status=%d len=%d", status, len(list))
	}

	status = s.Do(http.MethodGet, "/api/habits/"+habit.ID+"/stats", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("stats on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var stats struct {
		CurrentStreak int `json:"current_streak"`
	}
	status = s.Do(http.MethodGet, "/api/habits/"+habit.ID+"/stats", token, nil, &stats)
	if status != http.StatusOK || stats.CurrentStreak != 1 {
		t.Fatalf("stats on pro plan: status=%d body=%+v", status, stats)
	}
}

package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestHabits_PanoramaGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var habit struct{ ID string }
	status := s.Do(http.MethodPost, "/api/habits", token, map[string]string{"name": "Meditar"}, &habit)
	if status != http.StatusCreated || habit.ID == "" {
		t.Fatalf("create habit: status=%d", status)
	}
	status = s.Do(http.MethodPost, "/api/habits/"+habit.ID+"/log", token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("check-in: expected 204, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/habits/panorama", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("panorama on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var panorama []struct {
		HabitID       string    `json:"habit_id"`
		CurrentStreak int       `json:"current_streak"`
		BestStreak    int       `json:"best_streak"`
		WeeklyRates   []float64 `json:"weekly_rates"`
		Heatmap       []struct {
			Date      string `json:"date"`
			Completed bool   `json:"completed"`
		} `json:"heatmap"`
	}
	status = s.Do(http.MethodGet, "/api/habits/panorama", token, nil, &panorama)
	if status != http.StatusOK || len(panorama) != 1 {
		t.Fatalf("panorama on pro plan: status=%d body=%+v", status, panorama)
	}
	p := panorama[0]
	if p.HabitID != habit.ID || p.CurrentStreak != 1 || p.BestStreak != 1 {
		t.Fatalf("unexpected panorama: %+v", p)
	}
	if len(p.Heatmap) != 90 {
		t.Fatalf("expected 90-day heatmap, got %d", len(p.Heatmap))
	}
	if len(p.WeeklyRates) != 8 {
		t.Fatalf("expected 8 weekly rates, got %d", len(p.WeeklyRates))
	}
	if p.WeeklyRates[len(p.WeeklyRates)-1] <= 0 {
		t.Fatalf("expected latest week to have a positive completion rate, got %+v", p.WeeklyRates)
	}
}

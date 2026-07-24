package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestWorkouts_LogSessionHistoryProgressGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	body := map[string]interface{}{
		"session_date": "2026-07-01",
		"sets": []map[string]interface{}{
			{"exercise_name": "Sentadilla", "set_number": 1, "reps": 10, "weight_kg": 60},
			{"exercise_name": "Sentadilla", "set_number": 2, "reps": 8, "weight_kg": 65},
		},
	}
	var session struct {
		ID string `json:"id"`
	}
	status := s.Do(http.MethodPost, "/api/workouts", token, body, &session)
	if status != http.StatusCreated || session.ID == "" {
		t.Fatalf("log session: status=%d body=%+v", status, session)
	}

	var list []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/workouts", token, nil, &list)
	if status != http.StatusOK || len(list) != 1 {
		t.Fatalf("list sessions: status=%d len=%d", status, len(list))
	}

	status = s.Do(http.MethodGet, "/api/workouts/progress?exercise=Sentadilla", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("progress on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var progress struct {
		PersonalRecordKg float64 `json:"personal_record_kg"`
	}
	status = s.Do(http.MethodGet, "/api/workouts/progress?exercise=Sentadilla", token, nil, &progress)
	if status != http.StatusOK || progress.PersonalRecordKg != 65 {
		t.Fatalf("progress on pro plan: status=%d body=%+v", status, progress)
	}
}

func TestWorkouts_GetSessionOwnershipGating(t *testing.T) {
	s := testutil.NewServer(t)
	aliceToken := s.RegisterAndLogin("alice-get@example.com", "password123")
	bobToken := s.RegisterAndLogin("bob-get@example.com", "password123")

	body := map[string]interface{}{
		"session_date": "2026-07-01",
		"sets": []map[string]interface{}{
			{"exercise_name": "Sentadilla", "set_number": 1, "reps": 10, "weight_kg": 60},
		},
	}
	var session struct {
		ID string `json:"id"`
	}
	status := s.Do(http.MethodPost, "/api/workouts", aliceToken, body, &session)
	if status != http.StatusCreated || session.ID == "" {
		t.Fatalf("log session: status=%d body=%+v", status, session)
	}

	status = s.Do(http.MethodGet, "/api/workouts/"+session.ID, aliceToken, nil, nil)
	if status != http.StatusOK {
		t.Fatalf("owner get: expected 200, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/workouts/"+session.ID, bobToken, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("non-owner get without coach access: expected 403, got %d", status)
	}
}

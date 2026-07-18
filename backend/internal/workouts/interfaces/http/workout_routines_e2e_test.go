package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestWorkouts_RoutinesAndPersonalRecordsGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	body := map[string]interface{}{
		"session_date": "2026-07-01",
		"sets": []map[string]interface{}{
			{"exercise_name": "Sentadilla", "set_number": 1, "reps": 10, "weight_kg": 60},
			{"exercise_name": "Press banca", "set_number": 1, "reps": 8, "weight_kg": 40},
		},
	}
	status := s.Do(http.MethodPost, "/api/workouts", token, body, nil)
	if status != http.StatusCreated {
		t.Fatalf("log session: expected 201, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/workouts/personal-records", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("personal-records on free plan: expected 403, got %d", status)
	}

	status = s.Do(http.MethodPost, "/api/workouts/routines", token, map[string]interface{}{
		"name": "Empuje", "exercises": []string{"Press banca", "Fondos"},
	}, nil)
	if status != http.StatusForbidden {
		t.Fatalf("create routine on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var prs []struct {
		ExerciseName     string  `json:"exercise_name"`
		PersonalRecordKg float64 `json:"personal_record_kg"`
	}
	status = s.Do(http.MethodGet, "/api/workouts/personal-records", token, nil, &prs)
	if status != http.StatusOK || len(prs) != 2 {
		t.Fatalf("personal-records on pro plan: status=%d body=%+v", status, prs)
	}

	var routine struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Exercises []string `json:"exercises"`
	}
	status = s.Do(http.MethodPost, "/api/workouts/routines", token, map[string]interface{}{
		"name": "Empuje", "exercises": []string{"Press banca", "Fondos"},
	}, &routine)
	if status != http.StatusCreated || routine.ID == "" || len(routine.Exercises) != 2 {
		t.Fatalf("create routine: status=%d body=%+v", status, routine)
	}

	var routines []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/workouts/routines", token, nil, &routines)
	if status != http.StatusOK || len(routines) != 1 {
		t.Fatalf("list routines: status=%d len=%d", status, len(routines))
	}

	status = s.Do(http.MethodDelete, "/api/workouts/routines/"+routine.ID, token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("delete routine: expected 204, got %d", status)
	}

	routines = nil
	s.Do(http.MethodGet, "/api/workouts/routines", token, nil, &routines)
	if len(routines) != 0 {
		t.Fatalf("expected empty routines after delete, got %+v", routines)
	}
}

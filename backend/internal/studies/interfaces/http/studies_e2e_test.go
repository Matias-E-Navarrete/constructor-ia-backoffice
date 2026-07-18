package http_test

import (
	"net/http"
	"testing"
	"time"

	"rimu/backend/internal/platform/testutil"
)

func TestStudies_SubjectsSessionsOverviewGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var subject struct{ ID string }
	status := s.Do(http.MethodPost, "/api/studies/subjects", token, map[string]string{"name": "Cálculo I"}, &subject)
	if status != http.StatusCreated || subject.ID == "" {
		t.Fatalf("create subject: status=%d", status)
	}

	var subjects []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/studies/subjects", token, nil, &subjects)
	if status != http.StatusOK || len(subjects) != 1 {
		t.Fatalf("list subjects: status=%d len=%d", status, len(subjects))
	}

	today := time.Now().Format("2006-01-02")
	status = s.Do(http.MethodPost, "/api/studies/sessions", token, map[string]interface{}{
		"subject_id": subject.ID, "session_date": today, "duration_minutes": 45, "topic": "Límites",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("log session: expected 201, got %d", status)
	}

	var sessions []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/studies/sessions?subjectId="+subject.ID, token, nil, &sessions)
	if status != http.StatusOK || len(sessions) != 1 {
		t.Fatalf("list sessions: status=%d len=%d", status, len(sessions))
	}

	// A different user cannot log sessions against someone else's subject.
	otherToken := s.RegisterAndLogin("mallory@example.com", "password123")
	status = s.Do(http.MethodPost, "/api/studies/sessions", otherToken, map[string]interface{}{
		"subject_id": subject.ID, "session_date": today, "duration_minutes": 10,
	}, nil)
	if status != http.StatusForbidden {
		t.Fatalf("cross-user log session: expected 403, got %d", status)
	}

	status = s.Do(http.MethodGet, "/api/studies/overview", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("overview on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var overview []struct {
		SubjectID     string `json:"subject_id"`
		TotalMinutes  int    `json:"total_minutes"`
		CurrentStreak int    `json:"current_streak"`
	}
	status = s.Do(http.MethodGet, "/api/studies/overview", token, nil, &overview)
	if status != http.StatusOK || len(overview) != 1 {
		t.Fatalf("overview on pro plan: status=%d body=%+v", status, overview)
	}
	if overview[0].SubjectID != subject.ID || overview[0].TotalMinutes != 45 || overview[0].CurrentStreak != 1 {
		t.Fatalf("unexpected overview: %+v", overview[0])
	}

	status = s.Do(http.MethodDelete, "/api/studies/subjects/"+subject.ID, token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("delete subject: expected 204, got %d", status)
	}
}

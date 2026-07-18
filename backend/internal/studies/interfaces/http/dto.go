package http

import (
	"time"

	"rimu/backend/internal/studies/domain"
)

type createSubjectRequest struct {
	Name     string  `json:"name"`
	NoteSlug *string `json:"note_slug,omitempty"`
}

type subjectResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	NoteSlug *string `json:"note_slug,omitempty"`
}

func toSubjectResponse(s domain.Subject) subjectResponse {
	return subjectResponse{ID: s.ID, Name: s.Name, NoteSlug: s.NoteSlug}
}

type logSessionRequest struct {
	SubjectID       string `json:"subject_id"`
	SessionDate     string `json:"session_date"` // YYYY-MM-DD
	DurationMinutes int    `json:"duration_minutes"`
	Topic           string `json:"topic,omitempty"`
}

type sessionResponse struct {
	ID              string `json:"id"`
	SubjectID       string `json:"subject_id"`
	SessionDate     string `json:"session_date"`
	DurationMinutes int    `json:"duration_minutes"`
	Topic           string `json:"topic"`
}

func toSessionResponse(s domain.StudySession) sessionResponse {
	return sessionResponse{
		ID:              s.ID,
		SubjectID:       s.SubjectID,
		SessionDate:     s.SessionDate.Format("2006-01-02"),
		DurationMinutes: s.DurationMinutes,
		Topic:           s.Topic,
	}
}

type overviewResponse struct {
	SubjectID     string `json:"subject_id"`
	Name          string `json:"name"`
	SessionCount  int    `json:"session_count"`
	TotalMinutes  int    `json:"total_minutes"`
	CurrentStreak int    `json:"current_streak"`
}

func toOverviewResponse(o domain.SubjectOverview) overviewResponse {
	return overviewResponse{
		SubjectID:     o.SubjectID,
		Name:          o.Name,
		SessionCount:  o.SessionCount,
		TotalMinutes:  o.TotalMinutes,
		CurrentStreak: o.CurrentStreak,
	}
}

func parseDate(s string, fallback time.Time) (time.Time, error) {
	if s == "" {
		return fallback, nil
	}
	return time.Parse("2006-01-02", s)
}

package http

import "rimu/backend/internal/habits/domain"

type createHabitRequest struct {
	Name     string  `json:"name"`
	NoteSlug *string `json:"note_slug,omitempty"`
}

type updateHabitRequest struct {
	Name     *string `json:"name,omitempty"`
	NoteSlug *string `json:"note_slug,omitempty"`
}

type checkInRequest struct {
	Date string `json:"date,omitempty"` // YYYY-MM-DD, defaults to today
}

type habitResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	NoteSlug *string `json:"note_slug,omitempty"`
}

func toHabitResponse(h domain.Habit) habitResponse {
	return habitResponse{ID: h.ID, Name: h.Name, NoteSlug: h.NoteSlug}
}

type statsResponse struct {
	TotalLogs      int     `json:"total_logs"`
	CompletedLogs  int     `json:"completed_logs"`
	CompletionRate float64 `json:"completion_rate"`
	CurrentStreak  int     `json:"current_streak"`
}

func toStatsResponse(s domain.Stats) statsResponse {
	return statsResponse{
		TotalLogs:      s.TotalLogs,
		CompletedLogs:  s.CompletedLogs,
		CompletionRate: s.CompletionRate,
		CurrentStreak:  s.CurrentStreak,
	}
}

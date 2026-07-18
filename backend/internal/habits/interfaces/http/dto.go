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

type heatmapDayResponse struct {
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

type panoramaHabitResponse struct {
	HabitID       string               `json:"habit_id"`
	Name          string               `json:"name"`
	CurrentStreak int                  `json:"current_streak"`
	BestStreak    int                  `json:"best_streak"`
	Heatmap       []heatmapDayResponse `json:"heatmap"`
	WeeklyRates   []float64            `json:"weekly_rates"`
}

func toPanoramaResponse(p domain.PanoramaHabit) panoramaHabitResponse {
	heatmap := make([]heatmapDayResponse, 0, len(p.Heatmap))
	for _, d := range p.Heatmap {
		heatmap = append(heatmap, heatmapDayResponse{Date: d.Date, Completed: d.Completed})
	}
	return panoramaHabitResponse{
		HabitID:       p.HabitID,
		Name:          p.Name,
		CurrentStreak: p.CurrentStreak,
		BestStreak:    p.BestStreak,
		Heatmap:       heatmap,
		WeeklyRates:   p.WeeklyRates,
	}
}

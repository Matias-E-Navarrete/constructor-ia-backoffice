package http

import (
	"time"

	"rimu/backend/internal/workouts/application"
	"rimu/backend/internal/workouts/domain"
)

type setDTO struct {
	ExerciseName string  `json:"exercise_name"`
	SetNumber    int     `json:"set_number"`
	Reps         int     `json:"reps"`
	WeightKg     float64 `json:"weight_kg"`
}

type logSessionRequest struct {
	SessionDate string   `json:"session_date"` // YYYY-MM-DD
	Notes       string   `json:"notes,omitempty"`
	NoteSlug    *string  `json:"note_slug,omitempty"`
	Sets        []setDTO `json:"sets"`
}

type sessionResponse struct {
	ID          string   `json:"id"`
	SessionDate string   `json:"session_date"`
	Notes       string   `json:"notes"`
	NoteSlug    *string  `json:"note_slug,omitempty"`
	Sets        []setDTO `json:"sets"`
}

func toSessionResponse(s domain.Session) sessionResponse {
	sets := make([]setDTO, 0, len(s.Sets))
	for _, set := range s.Sets {
		sets = append(sets, setDTO{
			ExerciseName: set.ExerciseName, SetNumber: set.SetNumber, Reps: set.Reps, WeightKg: set.WeightKg,
		})
	}
	return sessionResponse{
		ID:          s.ID,
		SessionDate: s.SessionDate.Format("2006-01-02"),
		Notes:       s.Notes,
		NoteSlug:    s.NoteSlug,
		Sets:        sets,
	}
}

type progressEntryResponse struct {
	SessionDate string  `json:"session_date"`
	SetNumber   int     `json:"set_number"`
	Reps        int     `json:"reps"`
	WeightKg    float64 `json:"weight_kg"`
}

type progressResponse struct {
	Entries          []progressEntryResponse `json:"entries"`
	PersonalRecordKg float64                 `json:"personal_record_kg"`
}

func toProgressResponse(p application.ExerciseProgress) progressResponse {
	entries := make([]progressEntryResponse, 0, len(p.Entries))
	for _, e := range p.Entries {
		entries = append(entries, progressEntryResponse{
			SessionDate: e.SessionDate.Format("2006-01-02"), SetNumber: e.SetNumber, Reps: e.Reps, WeightKg: e.WeightKg,
		})
	}
	return progressResponse{Entries: entries, PersonalRecordKg: p.PersonalRecordKg}
}

func parseDate(s string, fallback time.Time) (time.Time, error) {
	if s == "" {
		return fallback, nil
	}
	return time.Parse("2006-01-02", s)
}

type createRoutineRequest struct {
	Name      string   `json:"name"`
	Exercises []string `json:"exercises"`
}

type routineResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Exercises []string `json:"exercises"`
}

func toRoutineResponse(r domain.Routine) routineResponse {
	return routineResponse{ID: r.ID, Name: r.Name, Exercises: r.Exercises}
}

type personalRecordResponse struct {
	ExerciseName     string  `json:"exercise_name"`
	PersonalRecordKg float64 `json:"personal_record_kg"`
}

func toPersonalRecordResponse(pr domain.ExercisePR) personalRecordResponse {
	return personalRecordResponse{ExerciseName: pr.ExerciseName, PersonalRecordKg: pr.PersonalRecordKg}
}

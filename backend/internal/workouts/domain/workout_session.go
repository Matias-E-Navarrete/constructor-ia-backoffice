package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("workouts: not found")

type Set struct {
	ExerciseName string
	SetNumber    int
	Reps         int
	WeightKg     float64
}

type Session struct {
	ID          string
	UserID      string
	SessionDate time.Time
	Notes       string
	NoteSlug    *string
	CreatedAt   time.Time
	Sets        []Set
}

// ExerciseEntry is one logged set for a given exercise, used to build a
// progress-over-time chart (the [PRO] "progress" endpoint).
type ExerciseEntry struct {
	SessionDate time.Time
	SetNumber   int
	Reps        int
	WeightKg    float64
}

// ExercisePR is one exercise's all-time personal record, used to render PR
// badges across every exercise the user has ever logged, without requiring
// them to look up each exercise name individually.
type ExercisePR struct {
	ExerciseName     string
	PersonalRecordKg float64
}

type Repository interface {
	Create(ctx context.Context, s *Session) error
	ListByUser(ctx context.Context, userID string) ([]Session, error)
	FindByID(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	ProgressForExercise(ctx context.Context, userID, exerciseName string) ([]ExerciseEntry, error)
	AllPersonalRecords(ctx context.Context, userID string) ([]ExercisePR, error)
}

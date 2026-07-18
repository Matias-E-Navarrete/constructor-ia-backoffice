package domain

import (
	"context"
	"time"
)

// Routine is a named, reusable list of exercises (a workout template) a user
// can start a session from, so they don't have to retype the same exercise
// list every time they train the same split.
type Routine struct {
	ID        string
	UserID    string
	Name      string
	Exercises []string
	CreatedAt time.Time
}

type RoutineRepository interface {
	CreateRoutine(ctx context.Context, r *Routine) error
	ListRoutinesByUser(ctx context.Context, userID string) ([]Routine, error)
	DeleteRoutine(ctx context.Context, userID, id string) error
}

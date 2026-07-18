package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

// GetPersonalRecords is a [PRO] use-case: the all-time PR for every exercise
// the user has ever logged, so the frontend can render PR badges without
// requiring the user to look up each exercise name one at a time.
type GetPersonalRecords struct {
	Repo domain.Repository
}

func (uc *GetPersonalRecords) Execute(ctx context.Context, userID string) ([]domain.ExercisePR, error) {
	return uc.Repo.AllPersonalRecords(ctx, userID)
}

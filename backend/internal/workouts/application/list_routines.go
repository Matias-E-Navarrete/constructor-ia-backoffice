package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

type ListRoutines struct {
	Repo domain.RoutineRepository
}

func (uc *ListRoutines) Execute(ctx context.Context, userID string) ([]domain.Routine, error) {
	return uc.Repo.ListRoutinesByUser(ctx, userID)
}

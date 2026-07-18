package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

type DeleteRoutine struct {
	Repo domain.RoutineRepository
}

func (uc *DeleteRoutine) Execute(ctx context.Context, userID, id string) error {
	return uc.Repo.DeleteRoutine(ctx, userID, id)
}

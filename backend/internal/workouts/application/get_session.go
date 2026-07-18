package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

type GetSession struct {
	Repo domain.Repository
}

func (uc *GetSession) Execute(ctx context.Context, sessionID string) (*domain.Session, error) {
	return uc.Repo.FindByID(ctx, sessionID)
}

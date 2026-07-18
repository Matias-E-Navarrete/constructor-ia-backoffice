package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

type DeleteSession struct {
	Repo domain.Repository
}

func (uc *DeleteSession) Execute(ctx context.Context, userID, sessionID string) error {
	s, err := uc.Repo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if s.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.Delete(ctx, sessionID)
}

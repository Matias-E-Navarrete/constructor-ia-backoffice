package application

import (
	"context"

	"rimu/backend/internal/habits/domain"
)

type DeleteHabit struct {
	Repo domain.Repository
}

func (uc *DeleteHabit) Execute(ctx context.Context, userID, habitID string) error {
	h, err := uc.Repo.FindByID(ctx, habitID)
	if err != nil {
		return err
	}
	if h.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.Delete(ctx, habitID)
}

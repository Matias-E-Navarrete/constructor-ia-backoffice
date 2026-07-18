package application

import (
	"context"

	"rimu/backend/internal/tasks/domain"
)

type DeleteTask struct {
	Repo domain.Repository
}

func (uc *DeleteTask) Execute(ctx context.Context, userID, taskID string) error {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if t.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.Delete(ctx, taskID)
}

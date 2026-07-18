package application

import (
	"context"
	"time"

	"rimu/backend/internal/tasks/domain"
)

type CompleteTask struct {
	Repo domain.Repository
}

func (uc *CompleteTask) Execute(ctx context.Context, userID, taskID string) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotOwner
	}
	t.Complete(time.Now())
	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

type UncompleteTask struct {
	Repo domain.Repository
}

func (uc *UncompleteTask) Execute(ctx context.Context, userID, taskID string) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotOwner
	}
	t.Uncomplete()
	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

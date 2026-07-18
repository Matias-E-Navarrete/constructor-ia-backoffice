package application

import (
	"context"
	"time"

	"rimu/backend/internal/tasks/domain"
)

// SetStatus moves a task between Kanban columns (todo/in_progress/done).
type SetStatus struct {
	Repo domain.Repository
}

func (uc *SetStatus) Execute(ctx context.Context, userID, taskID string, status domain.Status) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotOwner
	}
	t.Status = status
	if status == domain.StatusDone {
		t.Complete(time.Now())
	}
	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// SetQuadrant classifies a task into an Eisenhower quadrant.
type SetQuadrant struct {
	Repo domain.Repository
}

func (uc *SetQuadrant) Execute(ctx context.Context, userID, taskID string, quadrant domain.Quadrant) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotOwner
	}
	t.Quadrant = &quadrant
	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

package application

import (
	"context"
	"time"

	"rimu/backend/internal/tasks/domain"
)

// ScheduleTask assigns (or clears) the planner time-block for a task.
type ScheduleTask struct {
	Repo domain.Repository
}

func (uc *ScheduleTask) Execute(ctx context.Context, userID, taskID string, at *time.Time) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotOwner
	}
	t.ScheduledAt = at
	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

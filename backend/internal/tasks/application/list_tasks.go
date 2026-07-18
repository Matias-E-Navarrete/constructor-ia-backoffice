package application

import (
	"context"

	"rimu/backend/internal/tasks/domain"
)

type ListTasks struct {
	Repo domain.Repository
}

func (uc *ListTasks) Execute(ctx context.Context, userID string, filter domain.ListFilter) ([]domain.Task, error) {
	return uc.Repo.ListByUser(ctx, userID, filter)
}

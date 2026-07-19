package application

import (
	"context"

	"rimu/backend/internal/assistant/domain"
)

type ListMessages struct {
	Repo domain.Repository
}

func (uc *ListMessages) Execute(ctx context.Context, userID string) ([]domain.Message, error) {
	return uc.Repo.ListMessages(ctx, userID)
}

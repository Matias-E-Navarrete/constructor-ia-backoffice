package application

import (
	"context"

	"rimu/backend/internal/assistant/domain"
)

type ClearConversation struct {
	Repo domain.Repository
}

func (uc *ClearConversation) Execute(ctx context.Context, userID string) error {
	return uc.Repo.ClearMessages(ctx, userID)
}

package application

import (
	"context"

	"rimu/backend/internal/groups/domain"
)

type ListMyGroups struct {
	Repo domain.Repository
}

func (uc *ListMyGroups) Execute(ctx context.Context, userID string) ([]domain.Membership, error) {
	return uc.Repo.ListForUser(ctx, userID)
}

package application

import (
	"context"

	groupsdomain "rimu/backend/internal/groups/domain"
	"rimu/backend/internal/finance/domain"
)

type ListTransactions struct {
	Repo   domain.Repository
	Groups groupsdomain.Repository
}

func (uc *ListTransactions) Execute(ctx context.Context, userID string, groupID *string) ([]domain.Transaction, error) {
	if groupID == nil {
		return uc.Repo.ListByUser(ctx, userID)
	}

	isMember, err := uc.Groups.IsMember(ctx, *groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotGroupMember
	}
	return uc.Repo.ListByGroup(ctx, *groupID)
}

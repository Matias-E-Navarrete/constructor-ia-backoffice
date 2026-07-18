package application

import (
	"context"

	"rimu/backend/internal/tasks/domain"
)

// ReorderTasks persists the manual drag order for a user's task list. The
// repository scopes every update to userID, so IDs that don't belong to the
// caller are silently ignored rather than reordered.
type ReorderTasks struct {
	Repo domain.Repository
}

func (uc *ReorderTasks) Execute(ctx context.Context, userID string, orderedIDs []string) error {
	return uc.Repo.ReorderSortOrders(ctx, userID, orderedIDs)
}

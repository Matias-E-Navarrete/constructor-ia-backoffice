package application

import (
	"context"
	"time"

	"rimu/backend/internal/finance/domain"
)

// GetUpcomingBills is a [PRO] use-case: recurring transactions projected
// into the given window (e.g. "this month's upcoming bills").
type GetUpcomingBills struct {
	Repo domain.Repository
}

func (uc *GetUpcomingBills) Execute(ctx context.Context, userID string, from, to time.Time) ([]domain.UpcomingBill, error) {
	txs, err := uc.Repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return domain.ProjectUpcoming(txs, from, to), nil
}

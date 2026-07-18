package application

import (
	"context"

	"rimu/backend/internal/finance/domain"
)

// GetFinanceSummary is a [PRO] use-case: balance, income/expenses totals and
// a per-category breakdown.
type GetFinanceSummary struct {
	Repo domain.Repository
}

func (uc *GetFinanceSummary) Execute(ctx context.Context, userID string) (domain.Summary, error) {
	txs, err := uc.Repo.ListByUser(ctx, userID)
	if err != nil {
		return domain.Summary{}, err
	}
	return domain.ComputeSummary(txs), nil
}

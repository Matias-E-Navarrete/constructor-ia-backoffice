package application

import (
	"context"

	"rimu/backend/internal/user/domain"
)

// UpgradePlan is a mock stand-in for a real payment webhook: it flips the
// caller's plan to pro directly so the freemium gating is demonstrable
// end-to-end without wiring up Stripe/MercadoPago.
type UpgradePlan struct {
	Repo domain.Repository
}

func (uc *UpgradePlan) Execute(ctx context.Context, userID string) (*domain.User, error) {
	return uc.Repo.UpdatePlan(ctx, userID, domain.PlanPro)
}

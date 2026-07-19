package application

import (
	"context"

	userdomain "rimu/backend/internal/user/domain"
)

// ApplyProUpgrade is called from the Stripe webhook once a checkout session
// completes — the real-payment counterpart to the mock /api/me/upgrade
// endpoint, both ultimately just flipping the same plan field.
type ApplyProUpgrade struct {
	Users userdomain.Repository
}

func (uc *ApplyProUpgrade) Execute(ctx context.Context, userID string) error {
	_, err := uc.Users.UpdatePlan(ctx, userID, userdomain.PlanPro)
	return err
}

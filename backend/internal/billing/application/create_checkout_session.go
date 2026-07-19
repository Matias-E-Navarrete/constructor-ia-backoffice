package application

import (
	"context"

	billingdomain "rimu/backend/internal/billing/domain"
)

// StripeCheckout is the boundary to the Stripe API — implemented by
// infrastructure/stripe_client.go. Narrowed to just what this use-case
// needs so it doesn't depend on the concrete client type.
type StripeCheckout interface {
	CreateCheckoutSession(ctx context.Context, clientReferenceID, customerEmail, successURL, cancelURL string) (string, error)
}

type CreateCheckoutSession struct {
	Stripe     StripeCheckout
	SuccessURL string
	CancelURL  string
}

var ErrNotConfigured = billingdomain.ErrNotConfigured

func (uc *CreateCheckoutSession) Execute(ctx context.Context, userID, userEmail string) (string, error) {
	return uc.Stripe.CreateCheckoutSession(ctx, userID, userEmail, uc.SuccessURL, uc.CancelURL)
}

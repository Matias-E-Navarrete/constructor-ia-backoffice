package domain

import "errors"

// ErrNotConfigured means STRIPE_SECRET_KEY (and friends) aren't set on this
// deployment, so real checkout is unavailable and the caller should fall
// back to the mock upgrade path.
var ErrNotConfigured = errors.New("billing: no Stripe credentials configured for this deployment")

// ErrInvalidSignature means a webhook request's Stripe-Signature header
// didn't match the payload, so it's rejected rather than trusted.
var ErrInvalidSignature = errors.New("billing: webhook signature verification failed")

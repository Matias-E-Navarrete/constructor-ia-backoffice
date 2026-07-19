package domain

// WebhookEvent is the sliver of a Stripe event this app actually reads.
type WebhookEvent struct {
	Type              string
	ClientReferenceID string
}

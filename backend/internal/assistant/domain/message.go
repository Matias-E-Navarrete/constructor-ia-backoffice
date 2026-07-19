package domain

import (
	"context"
	"errors"
	"time"
)

// ErrNotConfigured means the deployment has no ANTHROPIC_API_KEY set, so the
// assistant feature is present but inert until an operator supplies one.
var ErrNotConfigured = errors.New("assistant: no API key configured for this deployment")

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is one turn in the user's single ongoing conversation with the
// assistant — kept as flat history per user rather than multiple named
// conversations, matching how a personal assistant is actually used.
type Message struct {
	ID        string
	UserID    string
	Role      Role
	Content   string
	CreatedAt time.Time
}

type Repository interface {
	SaveMessage(ctx context.Context, m *Message) error
	ListMessages(ctx context.Context, userID string) ([]Message, error)
	ClearMessages(ctx context.Context, userID string) error
}

// Completer is the boundary to whichever LLM provider answers the user —
// implemented by infrastructure/anthropic_client.go. Kept as a narrow
// interface so the application layer doesn't depend on any SDK.
type Completer interface {
	// Complete returns the assistant's reply given the full message history
	// (oldest first, ending with the latest user message).
	Complete(ctx context.Context, history []Message) (string, error)
}

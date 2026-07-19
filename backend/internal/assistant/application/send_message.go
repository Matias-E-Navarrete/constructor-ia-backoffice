package application

import (
	"context"

	"rimu/backend/internal/assistant/domain"
)

const maxHistoryMessages = 20

// SendMessage appends the user's message, asks the configured LLM for a
// reply grounded in the app's modules, and appends+returns that reply.
type SendMessage struct {
	Repo      domain.Repository
	Completer domain.Completer
}

func (uc *SendMessage) Execute(ctx context.Context, userID, content string) (*domain.Message, error) {
	userMsg := &domain.Message{UserID: userID, Role: domain.RoleUser, Content: content}
	if err := uc.Repo.SaveMessage(ctx, userMsg); err != nil {
		return nil, err
	}

	history, err := uc.Repo.ListMessages(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(history) > maxHistoryMessages {
		history = history[len(history)-maxHistoryMessages:]
	}

	reply, err := uc.Completer.Complete(ctx, history)
	if err != nil {
		return nil, err
	}

	assistantMsg := &domain.Message{UserID: userID, Role: domain.RoleAssistant, Content: reply}
	if err := uc.Repo.SaveMessage(ctx, assistantMsg); err != nil {
		return nil, err
	}
	return assistantMsg, nil
}

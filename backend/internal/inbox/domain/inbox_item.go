package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("inbox: not found")

type Item struct {
	ID        string
	UserID    string
	Content   string
	Pinned    bool
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, item *Item) error
	ListByUser(ctx context.Context, userID string) ([]Item, error)
	FindByID(ctx context.Context, id string) (*Item, error)
	SetPinned(ctx context.Context, id string, pinned bool) error
	Delete(ctx context.Context, id string) error
}

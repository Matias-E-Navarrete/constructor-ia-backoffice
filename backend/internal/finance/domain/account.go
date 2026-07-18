package domain

import (
	"context"
	"errors"
	"time"
)

var ErrAccountNotFound = errors.New("finance: account not found")

type Account struct {
	ID             string
	UserID         string
	Name           string
	Type           string
	Currency       string
	InitialBalance float64
	Active         bool
	CreatedAt      time.Time
}

type AccountRepository interface {
	Create(ctx context.Context, a *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	ListByUser(ctx context.Context, userID string) ([]Account, error)
	SetActive(ctx context.Context, id string, active bool) error
}

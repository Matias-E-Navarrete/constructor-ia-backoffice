package domain

import "context"

type Flag struct {
	Key         string
	Enabled     bool
	Description string
}

type Repository interface {
	List(ctx context.Context) ([]Flag, error)
	SetEnabled(ctx context.Context, key string, enabled bool) (Flag, error)
}

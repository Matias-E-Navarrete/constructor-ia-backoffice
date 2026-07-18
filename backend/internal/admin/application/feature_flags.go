package application

import (
	"context"

	ffdomain "rimu/backend/internal/platform/featureflags/domain"
)

// FlagRegistry is satisfied by featureflags/application.Registry.
type FlagRegistry interface {
	List(ctx context.Context) ([]ffdomain.Flag, error)
	Toggle(ctx context.Context, key string, enabled bool) (ffdomain.Flag, error)
}

type ListFeatureFlags struct {
	Registry FlagRegistry
}

func (uc *ListFeatureFlags) Execute(ctx context.Context) ([]ffdomain.Flag, error) {
	return uc.Registry.List(ctx)
}

type ToggleFeatureFlag struct {
	Registry FlagRegistry
}

func (uc *ToggleFeatureFlag) Execute(ctx context.Context, key string, enabled bool) (ffdomain.Flag, error) {
	return uc.Registry.Toggle(ctx, key, enabled)
}

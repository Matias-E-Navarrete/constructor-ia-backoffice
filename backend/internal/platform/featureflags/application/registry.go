// Package application holds the in-memory feature flag registry used by
// every module's middleware and by the admin toggle endpoint.
package application

import (
	"context"
	"sync"

	"rimu/backend/internal/platform/featureflags/domain"
)

// Registry caches flag state in memory so `IsEnabled` is cheap on every
// request. This is a single-instance MVP: a toggle updates Postgres and the
// in-memory map in the same call, so there is no cross-instance propagation
// to worry about yet.
type Registry struct {
	repo  domain.Repository
	mu    sync.RWMutex
	state map[string]bool
}

func NewRegistry(repo domain.Repository) *Registry {
	return &Registry{repo: repo, state: make(map[string]bool)}
}

func (r *Registry) Load(ctx context.Context) error {
	flags, err := r.repo.List(ctx)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, f := range flags {
		r.state[f.Key] = f.Enabled
	}
	return nil
}

// IsEnabled reports whether a feature is on. Unknown keys default to enabled
// so a missing seed row never accidentally locks a module out.
func (r *Registry) IsEnabled(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enabled, ok := r.state[key]
	if !ok {
		return true
	}
	return enabled
}

func (r *Registry) List(ctx context.Context) ([]domain.Flag, error) {
	return r.repo.List(ctx)
}

func (r *Registry) Toggle(ctx context.Context, key string, enabled bool) (domain.Flag, error) {
	flag, err := r.repo.SetEnabled(ctx, key, enabled)
	if err != nil {
		return domain.Flag{}, err
	}
	r.mu.Lock()
	r.state[flag.Key] = flag.Enabled
	r.mu.Unlock()
	return flag, nil
}

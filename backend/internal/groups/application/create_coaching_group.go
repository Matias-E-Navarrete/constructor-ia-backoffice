package application

import (
	"context"
	"errors"

	"rimu/backend/internal/groups/domain"
)

var ErrProRequired = errors.New("groups: coaching groups require the pro plan")

// CreateCoachingGroup creates a 'coaching' group with its creator as coach.
// Business rule: only pro accounts can run a coaching practice.
type CreateCoachingGroup struct {
	Repo domain.Repository
}

func (uc *CreateCoachingGroup) Execute(ctx context.Context, ownerUserID string, ownerIsPro bool, name string) (*domain.Group, error) {
	if !ownerIsPro {
		return nil, ErrProRequired
	}
	g := &domain.Group{Name: name, Kind: domain.KindCoaching, OwnerUserID: ownerUserID}
	if err := uc.Repo.Create(ctx, g, domain.RoleCoach); err != nil {
		return nil, err
	}
	return g, nil
}

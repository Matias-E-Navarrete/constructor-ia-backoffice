package application

import (
	"context"

	"rimu/backend/internal/groups/domain"
)

// CreateFamilyGroup creates a 'family' group. Free for every plan.
type CreateFamilyGroup struct {
	Repo domain.Repository
}

func (uc *CreateFamilyGroup) Execute(ctx context.Context, ownerUserID, name string) (*domain.Group, error) {
	g := &domain.Group{Name: name, Kind: domain.KindFamily, OwnerUserID: ownerUserID}
	if err := uc.Repo.Create(ctx, g, domain.RoleOwner); err != nil {
		return nil, err
	}
	return g, nil
}

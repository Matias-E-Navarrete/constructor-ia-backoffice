package application

import (
	"context"

	userdomain "rimu/backend/internal/user/domain"
)

type ListUsers struct {
	Repo userdomain.Repository
}

func (uc *ListUsers) Execute(ctx context.Context) ([]*userdomain.User, error) {
	return uc.Repo.List(ctx)
}

type SetUserPlan struct {
	Repo userdomain.Repository
}

func (uc *SetUserPlan) Execute(ctx context.Context, userID string, plan userdomain.Plan) (*userdomain.User, error) {
	return uc.Repo.UpdatePlan(ctx, userID, plan)
}

type SetUserRole struct {
	Repo userdomain.Repository
}

func (uc *SetUserRole) Execute(ctx context.Context, userID string, role userdomain.Role) (*userdomain.User, error) {
	return uc.Repo.UpdateRole(ctx, userID, role)
}

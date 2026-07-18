package application

import (
	"context"
	"errors"

	"rimu/backend/internal/user/domain"
)

var ErrInvalidPlannerHours = errors.New("user: planner hours must satisfy 0 <= start < end <= 24")

type UpdatePlannerHours struct {
	Repo domain.Repository
}

func (uc *UpdatePlannerHours) Execute(ctx context.Context, userID string, startHour, endHour int) (*domain.User, error) {
	if startHour < 0 || endHour > 24 || startHour >= endHour {
		return nil, ErrInvalidPlannerHours
	}
	return uc.Repo.UpdatePlannerHours(ctx, userID, startHour, endHour)
}

package application

import (
	"context"

	groupsdomain "rimu/backend/internal/groups/domain"
	"rimu/backend/internal/workouts/domain"
)

type GetSession struct {
	Repo   domain.Repository
	Groups groupsdomain.Repository
}

func (uc *GetSession) Execute(ctx context.Context, requesterUserID, sessionID string) (*domain.Session, error) {
	s, err := uc.Repo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if s.UserID == requesterUserID {
		return s, nil
	}

	allowed, err := uc.Groups.HasCoachAccess(ctx, requesterUserID, s.UserID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrNoCoachAccess
	}
	return s, nil
}

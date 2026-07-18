package application

import (
	"context"
	"errors"

	groupsdomain "rimu/backend/internal/groups/domain"
	"rimu/backend/internal/workouts/domain"
)

var ErrNoCoachAccess = errors.New("workouts: you do not have coach access to this user")
var ErrNotOwner = errors.New("workouts: you do not own this session")

// ListSessions mirrors habits.ListHabits: defaults to the caller, allows a
// coach to pass a client's userId if an active coaching relationship exists.
type ListSessions struct {
	Repo   domain.Repository
	Groups groupsdomain.Repository
}

func (uc *ListSessions) Execute(ctx context.Context, requesterUserID, targetUserID string) ([]domain.Session, error) {
	if targetUserID == "" || targetUserID == requesterUserID {
		return uc.Repo.ListByUser(ctx, requesterUserID)
	}

	allowed, err := uc.Groups.HasCoachAccess(ctx, requesterUserID, targetUserID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrNoCoachAccess
	}
	return uc.Repo.ListByUser(ctx, targetUserID)
}

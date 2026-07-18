package application

import (
	"context"
	"errors"

	"rimu/backend/internal/habits/domain"
	groupsdomain "rimu/backend/internal/groups/domain"
)

var ErrNoCoachAccess = errors.New("habits: you do not have coach access to this user")

// ListHabits defaults to the caller's own habits. Passing a different
// targetUserID (a coach viewing a client) requires an active coaching
// relationship, checked via the groups bounded context.
type ListHabits struct {
	Repo   domain.Repository
	Groups groupsdomain.Repository
}

func (uc *ListHabits) Execute(ctx context.Context, requesterUserID, targetUserID string) ([]domain.Habit, error) {
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

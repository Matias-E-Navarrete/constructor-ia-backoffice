package application

import (
	"context"
	"time"

	"rimu/backend/internal/habits/domain"
)

// GetHabitStats is a [PRO] use-case: completion rate + current streak.
type GetHabitStats struct {
	Repo domain.Repository
}

func (uc *GetHabitStats) Execute(ctx context.Context, userID, habitID string) (domain.Stats, error) {
	h, err := uc.Repo.FindByID(ctx, habitID)
	if err != nil {
		return domain.Stats{}, err
	}
	if h.UserID != userID {
		return domain.Stats{}, ErrNotOwner
	}

	logs, err := uc.Repo.ListLogs(ctx, habitID)
	if err != nil {
		return domain.Stats{}, err
	}
	return domain.ComputeStats(logs, time.Now()), nil
}

package application

import (
	"context"
	"errors"
	"time"

	"rimu/backend/internal/habits/domain"
)

var ErrNotOwner = errors.New("habits: you do not own this habit")

type CheckInHabit struct {
	Repo domain.Repository
}

func (uc *CheckInHabit) Execute(ctx context.Context, userID, habitID string, date time.Time) error {
	h, err := uc.Repo.FindByID(ctx, habitID)
	if err != nil {
		return err
	}
	if h.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.UpsertLog(ctx, h.CheckIn(date))
}

package application

import (
	"context"

	"rimu/backend/internal/habits/domain"
)

type UpdateHabitInput struct {
	UserID   string
	HabitID  string
	Name     *string
	NoteSlug *string
}

type UpdateHabit struct {
	Repo domain.Repository
}

func (uc *UpdateHabit) Execute(ctx context.Context, in UpdateHabitInput) (*domain.Habit, error) {
	h, err := uc.Repo.FindByID(ctx, in.HabitID)
	if err != nil {
		return nil, err
	}
	if h.UserID != in.UserID {
		return nil, ErrNotOwner
	}
	if in.Name != nil {
		h.Name = *in.Name
	}
	if in.NoteSlug != nil {
		h.NoteSlug = in.NoteSlug
	}
	if err := uc.Repo.Update(ctx, h); err != nil {
		return nil, err
	}
	return h, nil
}

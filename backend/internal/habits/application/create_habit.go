package application

import (
	"context"

	"rimu/backend/internal/habits/domain"
)

type CreateHabit struct {
	Repo domain.Repository
}

type CreateHabitInput struct {
	UserID   string
	Name     string
	NoteSlug *string
}

func (uc *CreateHabit) Execute(ctx context.Context, in CreateHabitInput) (*domain.Habit, error) {
	h := &domain.Habit{UserID: in.UserID, Name: in.Name, NoteSlug: in.NoteSlug}
	if err := uc.Repo.Create(ctx, h); err != nil {
		return nil, err
	}
	return h, nil
}

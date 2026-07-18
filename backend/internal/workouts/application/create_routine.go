package application

import (
	"context"
	"errors"

	"rimu/backend/internal/workouts/domain"
)

var ErrRoutineNeedsExercises = errors.New("workouts: a routine needs at least one exercise")

type CreateRoutine struct {
	Repo domain.RoutineRepository
}

type CreateRoutineInput struct {
	UserID    string
	Name      string
	Exercises []string
}

func (uc *CreateRoutine) Execute(ctx context.Context, in CreateRoutineInput) (*domain.Routine, error) {
	if len(in.Exercises) == 0 {
		return nil, ErrRoutineNeedsExercises
	}
	r := &domain.Routine{UserID: in.UserID, Name: in.Name, Exercises: in.Exercises}
	if err := uc.Repo.CreateRoutine(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

package application

import (
	"context"

	"rimu/backend/internal/workouts/domain"
)

// GetExerciseProgress is a [PRO] use-case: the weight/reps time series for a
// given exercise, plus the personal record (max weight lifted).
type GetExerciseProgress struct {
	Repo domain.Repository
}

type ExerciseProgress struct {
	Entries          []domain.ExerciseEntry
	PersonalRecordKg float64
}

func (uc *GetExerciseProgress) Execute(ctx context.Context, userID, exerciseName string) (ExerciseProgress, error) {
	entries, err := uc.Repo.ProgressForExercise(ctx, userID, exerciseName)
	if err != nil {
		return ExerciseProgress{}, err
	}

	var pr float64
	for _, e := range entries {
		if e.WeightKg > pr {
			pr = e.WeightKg
		}
	}
	return ExerciseProgress{Entries: entries, PersonalRecordKg: pr}, nil
}

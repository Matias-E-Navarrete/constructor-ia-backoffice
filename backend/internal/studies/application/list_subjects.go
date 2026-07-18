package application

import (
	"context"

	"rimu/backend/internal/studies/domain"
)

type ListSubjects struct {
	Repo domain.Repository
}

func (uc *ListSubjects) Execute(ctx context.Context, userID string) ([]domain.Subject, error) {
	return uc.Repo.ListSubjectsByUser(ctx, userID)
}

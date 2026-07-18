package application

import (
	"context"

	"rimu/backend/internal/studies/domain"
)

type DeleteSubject struct {
	Repo domain.Repository
}

func (uc *DeleteSubject) Execute(ctx context.Context, userID, id string) error {
	return uc.Repo.DeleteSubject(ctx, userID, id)
}

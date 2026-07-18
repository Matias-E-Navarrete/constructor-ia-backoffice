package application

import (
	"context"

	"rimu/backend/internal/notes/domain"
)

type ListNotes struct {
	Repo domain.Repository
}

func (uc *ListNotes) Execute(ctx context.Context, userID string) ([]domain.Note, error) {
	return uc.Repo.ListByUser(ctx, userID)
}

type GetNote struct {
	Repo domain.Repository
}

func (uc *GetNote) Execute(ctx context.Context, userID, slug string) (*domain.Note, error) {
	return uc.Repo.FindBySlug(ctx, userID, slug)
}

type DeleteNote struct {
	Repo domain.Repository
}

func (uc *DeleteNote) Execute(ctx context.Context, userID, slug string) error {
	return uc.Repo.Delete(ctx, userID, slug)
}

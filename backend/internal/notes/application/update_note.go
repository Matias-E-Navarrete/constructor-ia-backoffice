package application

import (
	"context"

	"rimu/backend/internal/notes/domain"
)

type UpdateNoteInput struct {
	UserID string
	Slug   string
	Title  string
	Body   string
}

type UpdateNote struct {
	Repo domain.Repository
}

func (uc *UpdateNote) Execute(ctx context.Context, in UpdateNoteInput) (*domain.Note, error) {
	n, err := uc.Repo.FindBySlug(ctx, in.UserID, in.Slug)
	if err != nil {
		return nil, err
	}
	n.Title = in.Title
	n.BodyMarkdown = in.Body
	if err := uc.Repo.Update(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

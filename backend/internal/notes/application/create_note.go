package application

import (
	"context"

	"rimu/backend/internal/notes/domain"
)

type CreateNoteInput struct {
	UserID string
	Title  string
	Body   string
	Kind   domain.Kind
}

type CreateNote struct {
	Repo domain.Repository
}

func (uc *CreateNote) Execute(ctx context.Context, in CreateNoteInput) (*domain.Note, error) {
	if in.Kind == "" {
		in.Kind = domain.KindFreeform
	}
	n := &domain.Note{
		UserID:       in.UserID,
		Slug:         domain.Slugify(in.Title),
		Title:        in.Title,
		BodyMarkdown: in.Body,
		Kind:         in.Kind,
	}
	if err := uc.Repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

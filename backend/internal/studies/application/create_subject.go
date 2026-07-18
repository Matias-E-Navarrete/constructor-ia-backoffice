package application

import (
	"context"

	"rimu/backend/internal/studies/domain"
)

type CreateSubject struct {
	Repo domain.Repository
}

type CreateSubjectInput struct {
	UserID   string
	Name     string
	NoteSlug *string
}

func (uc *CreateSubject) Execute(ctx context.Context, in CreateSubjectInput) (*domain.Subject, error) {
	s := &domain.Subject{UserID: in.UserID, Name: in.Name, NoteSlug: in.NoteSlug}
	if err := uc.Repo.CreateSubject(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

package application

import (
	"context"
	"time"

	"rimu/backend/internal/workouts/domain"
)

type LogSessionInput struct {
	UserID      string
	SessionDate time.Time
	Notes       string
	NoteSlug    *string
	Sets        []domain.Set
}

type LogSession struct {
	Repo domain.Repository
}

func (uc *LogSession) Execute(ctx context.Context, in LogSessionInput) (*domain.Session, error) {
	s := &domain.Session{
		UserID:      in.UserID,
		SessionDate: in.SessionDate,
		Notes:       in.Notes,
		NoteSlug:    in.NoteSlug,
		Sets:        in.Sets,
	}
	if err := uc.Repo.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

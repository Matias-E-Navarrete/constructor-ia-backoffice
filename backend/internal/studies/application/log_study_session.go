package application

import (
	"context"
	"time"

	"rimu/backend/internal/studies/domain"
)

type LogStudySession struct {
	Repo domain.Repository
}

type LogStudySessionInput struct {
	UserID          string
	SubjectID       string
	SessionDate     time.Time
	DurationMinutes int
	Topic           string
}

func (uc *LogStudySession) Execute(ctx context.Context, in LogStudySessionInput) (*domain.StudySession, error) {
	subject, err := uc.Repo.FindSubjectByID(ctx, in.SubjectID)
	if err != nil {
		return nil, err
	}
	if subject.UserID != in.UserID {
		return nil, domain.ErrNotOwner
	}

	session := &domain.StudySession{
		UserID:          in.UserID,
		SubjectID:       in.SubjectID,
		SessionDate:     in.SessionDate,
		DurationMinutes: in.DurationMinutes,
		Topic:           in.Topic,
	}
	if err := uc.Repo.LogSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

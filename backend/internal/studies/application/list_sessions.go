package application

import (
	"context"

	"rimu/backend/internal/studies/domain"
)

// ListSessions defaults to every session across every subject; passing a
// subjectID scopes the list to that one subject.
type ListSessions struct {
	Repo domain.Repository
}

func (uc *ListSessions) Execute(ctx context.Context, userID, subjectID string) ([]domain.StudySession, error) {
	if subjectID == "" {
		return uc.Repo.ListSessionsByUser(ctx, userID)
	}

	subject, err := uc.Repo.FindSubjectByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if subject.UserID != userID {
		return nil, domain.ErrNotOwner
	}
	return uc.Repo.ListSessionsBySubject(ctx, subjectID)
}

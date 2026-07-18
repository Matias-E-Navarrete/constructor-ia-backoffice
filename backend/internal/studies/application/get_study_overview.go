package application

import (
	"context"
	"time"

	"rimu/backend/internal/studies/domain"
)

// GetStudyOverview is a [PRO] use-case: total time invested and current
// streak per subject.
type GetStudyOverview struct {
	Repo domain.Repository
}

func (uc *GetStudyOverview) Execute(ctx context.Context, userID string) ([]domain.SubjectOverview, error) {
	subjects, err := uc.Repo.ListSubjectsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	overview := make([]domain.SubjectOverview, 0, len(subjects))
	for _, s := range subjects {
		sessions, err := uc.Repo.ListSessionsBySubject(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		overview = append(overview, domain.ComputeOverview(s, sessions, now))
	}
	return overview, nil
}

package domain

import "time"

// SubjectOverview is the [PRO] Estudios dashboard's per-subject summary:
// total time invested, how many sessions, and the current daily study streak.
type SubjectOverview struct {
	SubjectID     string
	Name          string
	SessionCount  int
	TotalMinutes  int
	CurrentStreak int
}

// ComputeOverview derives totals and the current streak from a subject's raw
// session rows. Kept as a pure function so it's unit-testable without a
// database, mirroring habits.ComputeStats/ComputePanorama.
func ComputeOverview(s Subject, sessions []StudySession, asOf time.Time) SubjectOverview {
	studiedOn := make(map[string]bool, len(sessions))
	totalMinutes := 0
	for _, sess := range sessions {
		studiedOn[sess.SessionDate.Format("2006-01-02")] = true
		totalMinutes += sess.DurationMinutes
	}

	streak := 0
	for d := asOf; ; d = d.AddDate(0, 0, -1) {
		if studiedOn[d.Format("2006-01-02")] {
			streak++
			continue
		}
		break
	}

	return SubjectOverview{
		SubjectID:     s.ID,
		Name:          s.Name,
		SessionCount:  len(sessions),
		TotalMinutes:  totalMinutes,
		CurrentStreak: streak,
	}
}

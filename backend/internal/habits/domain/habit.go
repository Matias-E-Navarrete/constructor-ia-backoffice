package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("habits: not found")

type Habit struct {
	ID         string
	UserID     string
	Name       string
	NoteSlug   *string
	CreatedAt  time.Time
	ArchivedAt *time.Time
}

type HabitLog struct {
	HabitID   string
	LogDate   time.Time
	Completed bool
}

// CheckIn produces the log entry for marking a habit done on a given day.
// The repository is responsible for upserting it (one row per habit per day).
func (h *Habit) CheckIn(date time.Time) HabitLog {
	return HabitLog{HabitID: h.ID, LogDate: date, Completed: true}
}

type Stats struct {
	TotalLogs      int
	CompletedLogs  int
	CompletionRate float64
	CurrentStreak  int
}

// ComputeStats derives completion rate and current streak from the raw log
// rows. Kept as a pure function so it's trivially unit-testable without a
// database.
func ComputeStats(logs []HabitLog, asOf time.Time) Stats {
	byDate := make(map[string]bool, len(logs))
	completed := 0
	for _, l := range logs {
		byDate[l.LogDate.Format("2006-01-02")] = l.Completed
		if l.Completed {
			completed++
		}
	}

	streak := 0
	for d := asOf; ; d = d.AddDate(0, 0, -1) {
		if done, ok := byDate[d.Format("2006-01-02")]; ok && done {
			streak++
			continue
		}
		break
	}

	rate := 0.0
	if len(logs) > 0 {
		rate = float64(completed) / float64(len(logs))
	}

	return Stats{
		TotalLogs:      len(logs),
		CompletedLogs:  completed,
		CompletionRate: rate,
		CurrentStreak:  streak,
	}
}

type Repository interface {
	Create(ctx context.Context, h *Habit) error
	ListByUser(ctx context.Context, userID string) ([]Habit, error)
	FindByID(ctx context.Context, id string) (*Habit, error)
	Update(ctx context.Context, h *Habit) error
	Delete(ctx context.Context, id string) error
	UpsertLog(ctx context.Context, log HabitLog) error
	ListLogs(ctx context.Context, habitID string) ([]HabitLog, error)
}

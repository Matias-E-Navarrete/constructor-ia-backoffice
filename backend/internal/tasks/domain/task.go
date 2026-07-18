package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("tasks: not found")

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Quadrant string

const (
	QuadrantDoNow    Quadrant = "do_now"
	QuadrantSchedule Quadrant = "schedule"
	QuadrantDelegate Quadrant = "delegate"
	QuadrantEliminate Quadrant = "eliminate"
)

type RepeatRule string

const (
	RepeatNone    RepeatRule = "none"
	RepeatDaily   RepeatRule = "daily"
	RepeatWeekly  RepeatRule = "weekly"
	RepeatMonthly RepeatRule = "monthly"
)

type Subtask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Task struct {
	ID           string
	UserID       string
	Title        string
	Description  string
	Category     string
	Subcategory  *string
	DueDate      *time.Time
	DueTime      *string
	RepeatRule   RepeatRule
	Priority     Priority
	Status       Status
	Quadrant     *Quadrant
	SortOrder    int
	ScheduledAt  *time.Time
	Subtasks     []Subtask
	NoteSlug     *string
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Task) Complete(now time.Time) {
	t.Status = StatusDone
	t.CompletedAt = &now
}

func (t *Task) Uncomplete() {
	t.Status = StatusTodo
	t.CompletedAt = nil
}

func (t *Task) IsOverdue(now time.Time) bool {
	if t.Status == StatusDone || t.DueDate == nil {
		return false
	}
	return t.DueDate.Before(now.Truncate(24 * time.Hour))
}

type ListFilter struct {
	Status   *Status
	Quadrant *Quadrant
	From     *time.Time
	To       *time.Time
}

type Repository interface {
	Create(ctx context.Context, t *Task) error
	Update(ctx context.Context, t *Task) error
	FindByID(ctx context.Context, id string) (*Task, error)
	ListByUser(ctx context.Context, userID string, filter ListFilter) ([]Task, error)
	Delete(ctx context.Context, id string) error
	ReorderSortOrders(ctx context.Context, userID string, orderedIDs []string) error
}

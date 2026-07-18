package application

import (
	"context"
	"time"

	"rimu/backend/internal/tasks/domain"
)

type CreateTaskInput struct {
	UserID      string
	Title       string
	Description string
	Category    string
	Subcategory *string
	DueDate     *time.Time
	DueTime     *string
	RepeatRule  domain.RepeatRule
	Priority    domain.Priority
	Subtasks    []domain.Subtask
	NoteSlug    *string
}

type CreateTask struct {
	Repo domain.Repository
}

func (uc *CreateTask) Execute(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	category := in.Category
	if category == "" {
		category = "General"
	}
	priority := in.Priority
	if priority == "" {
		priority = domain.PriorityNormal
	}
	repeat := in.RepeatRule
	if repeat == "" {
		repeat = domain.RepeatNone
	}

	t := &domain.Task{
		UserID:      in.UserID,
		Title:       in.Title,
		Description: in.Description,
		Category:    category,
		Subcategory: in.Subcategory,
		DueDate:     in.DueDate,
		DueTime:     in.DueTime,
		RepeatRule:  repeat,
		Priority:    priority,
		Status:      domain.StatusTodo,
		Subtasks:    in.Subtasks,
		NoteSlug:    in.NoteSlug,
	}
	if t.Subtasks == nil {
		t.Subtasks = []domain.Subtask{}
	}
	if err := uc.Repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

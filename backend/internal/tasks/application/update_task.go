package application

import (
	"context"
	"errors"
	"time"

	"rimu/backend/internal/tasks/domain"
)

var ErrNotOwner = errors.New("tasks: you do not own this task")

type UpdateTaskInput struct {
	UserID      string
	TaskID      string
	Title       *string
	Description *string
	Category    *string
	Subcategory *string
	DueDate     *time.Time
	ClearDate   bool
	DueTime     *string
	RepeatRule  *domain.RepeatRule
	Priority    *domain.Priority
}

type UpdateTask struct {
	Repo domain.Repository
}

func (uc *UpdateTask) Execute(ctx context.Context, in UpdateTaskInput) (*domain.Task, error) {
	t, err := uc.Repo.FindByID(ctx, in.TaskID)
	if err != nil {
		return nil, err
	}
	if t.UserID != in.UserID {
		return nil, ErrNotOwner
	}

	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.Category != nil {
		t.Category = *in.Category
	}
	if in.Subcategory != nil {
		t.Subcategory = in.Subcategory
	}
	if in.ClearDate {
		t.DueDate = nil
	} else if in.DueDate != nil {
		t.DueDate = in.DueDate
	}
	if in.DueTime != nil {
		t.DueTime = in.DueTime
	}
	if in.RepeatRule != nil {
		t.RepeatRule = *in.RepeatRule
	}
	if in.Priority != nil {
		t.Priority = *in.Priority
	}

	if err := uc.Repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

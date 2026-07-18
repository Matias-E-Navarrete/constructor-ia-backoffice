package http

import (
	"time"

	"rimu/backend/internal/tasks/domain"
)

type subtaskDTO struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type createTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Category    string       `json:"category,omitempty"`
	Subcategory *string      `json:"subcategory,omitempty"`
	DueDate     *string      `json:"due_date,omitempty"` // YYYY-MM-DD
	DueTime     *string      `json:"due_time,omitempty"` // HH:MM
	RepeatRule  string       `json:"repeat_rule,omitempty"`
	Priority    string       `json:"priority,omitempty"`
	Subtasks    []subtaskDTO `json:"subtasks,omitempty"`
	NoteSlug    *string      `json:"note_slug,omitempty"`
}

type updateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Subcategory *string `json:"subcategory,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	ClearDate   bool    `json:"clear_date,omitempty"`
	DueTime     *string `json:"due_time,omitempty"`
	RepeatRule  *string `json:"repeat_rule,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

type setStatusRequest struct {
	Status string `json:"status"`
}

type setQuadrantRequest struct {
	Quadrant string `json:"quadrant"`
}

type scheduleRequest struct {
	ScheduledAt *string `json:"scheduled_at,omitempty"` // RFC3339, null clears
}

type reorderRequest struct {
	OrderedIDs []string `json:"ordered_ids"`
}

type taskResponse struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Subcategory *string      `json:"subcategory,omitempty"`
	DueDate     *string      `json:"due_date,omitempty"`
	DueTime     *string      `json:"due_time,omitempty"`
	RepeatRule  string       `json:"repeat_rule"`
	Priority    string       `json:"priority"`
	Status      string       `json:"status"`
	Quadrant    *string      `json:"quadrant,omitempty"`
	SortOrder   int          `json:"sort_order"`
	ScheduledAt *string      `json:"scheduled_at,omitempty"`
	Subtasks    []subtaskDTO `json:"subtasks"`
	NoteSlug    *string      `json:"note_slug,omitempty"`
	CompletedAt *string      `json:"completed_at,omitempty"`
	Overdue     bool         `json:"overdue"`
}

func toTaskResponse(t domain.Task) taskResponse {
	resp := taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Category:    t.Category,
		Subcategory: t.Subcategory,
		DueTime:     t.DueTime,
		RepeatRule:  string(t.RepeatRule),
		Priority:    string(t.Priority),
		Status:      string(t.Status),
		SortOrder:   t.SortOrder,
		NoteSlug:    t.NoteSlug,
		Overdue:     t.IsOverdue(time.Now()),
	}
	resp.Subtasks = make([]subtaskDTO, 0, len(t.Subtasks))
	for _, s := range t.Subtasks {
		resp.Subtasks = append(resp.Subtasks, subtaskDTO{Title: s.Title, Done: s.Done})
	}
	if t.DueDate != nil {
		s := t.DueDate.Format("2006-01-02")
		resp.DueDate = &s
	}
	if t.Quadrant != nil {
		s := string(*t.Quadrant)
		resp.Quadrant = &s
	}
	if t.ScheduledAt != nil {
		s := t.ScheduledAt.Format(time.RFC3339)
		resp.ScheduledAt = &s
	}
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &s
	}
	return resp
}

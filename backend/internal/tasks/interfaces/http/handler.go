package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
	"rimu/backend/internal/tasks/application"
	"rimu/backend/internal/tasks/domain"
)

type Handler struct {
	Create      *application.CreateTask
	Update      *application.UpdateTask
	Delete      *application.DeleteTask
	List        *application.ListTasks
	Complete    *application.CompleteTask
	Uncomplete  *application.UncompleteTask
	SetStatus   *application.SetStatus
	SetQuadrant *application.SetQuadrant
	Reorder     *application.ReorderTasks
	Schedule    *application.ScheduleTask
}

func (h *Handler) HandleCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createTaskRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	subtasks := make([]domain.Subtask, 0, len(req.Subtasks))
	for _, s := range req.Subtasks {
		subtasks = append(subtasks, domain.Subtask{Title: s.Title, Done: s.Done})
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "due_date must be YYYY-MM-DD")
			return
		}
		dueDate = &parsed
	}

	task, err := h.Create.Execute(r.Context(), application.CreateTaskInput{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Subcategory: req.Subcategory,
		DueDate:     dueDate,
		DueTime:     req.DueTime,
		RepeatRule:  domain.RepeatRule(req.RepeatRule),
		Priority:    domain.Priority(req.Priority),
		Subtasks:    subtasks,
		NoteSlug:    req.NoteSlug,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toTaskResponse(*task))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	filter := domain.ListFilter{}
	q := r.URL.Query()
	if status := q.Get("status"); status != "" {
		s := domain.Status(status)
		filter.Status = &s
	}
	if quadrant := q.Get("quadrant"); quadrant != "" {
		qd := domain.Quadrant(quadrant)
		filter.Quadrant = &qd
	}
	if from := q.Get("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.From = &t
		}
	}
	if to := q.Get("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			filter.To = &t
		}
	}

	tasks, err := h.List.Execute(r.Context(), user.ID, filter)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, toTaskResponse(t))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleUpdate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req updateTaskRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	in := application.UpdateTaskInput{
		UserID: user.ID, TaskID: id,
		Title: req.Title, Description: req.Description, Category: req.Category, Subcategory: req.Subcategory,
		DueTime: req.DueTime, ClearDate: req.ClearDate,
	}
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "due_date must be YYYY-MM-DD")
			return
		}
		in.DueDate = &parsed
	}
	if req.RepeatRule != nil {
		rr := domain.RepeatRule(*req.RepeatRule)
		in.RepeatRule = &rr
	}
	if req.Priority != nil {
		p := domain.Priority(*req.Priority)
		in.Priority = &p
	}

	task, err := h.Update.Execute(r.Context(), in)
	if err == application.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleDelete(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	err := h.Delete.Execute(r.Context(), user.ID, id)
	if err == application.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *Handler) HandleComplete(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	task, err := h.Complete.Execute(r.Context(), user.ID, id)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleUncomplete(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	task, err := h.Uncomplete.Execute(r.Context(), user.ID, id)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleSetStatus(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req setStatusRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	task, err := h.SetStatus.Execute(r.Context(), user.ID, id, domain.Status(req.Status))
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleSetQuadrant(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req setQuadrantRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	task, err := h.SetQuadrant.Execute(r.Context(), user.ID, id, domain.Quadrant(req.Quadrant))
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleSchedule(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req scheduleRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	var at *time.Time
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "scheduled_at must be RFC3339")
			return
		}
		at = &parsed
	}

	task, err := h.Schedule.Execute(r.Context(), user.ID, id, at)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toTaskResponse(*task))
}

func (h *Handler) HandleReorder(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req reorderRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	if err := h.Reorder.Execute(r.Context(), user.ID, req.OrderedIDs); err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/habits/application"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Create    *application.CreateHabit
	Update    *application.UpdateHabit
	Delete    *application.DeleteHabit
	CheckIn   *application.CheckInHabit
	List      *application.ListHabits
	GetStats  *application.GetHabitStats
}

func (h *Handler) HandleCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createHabitRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	habit, err := h.Create.Execute(r.Context(), application.CreateHabitInput{
		UserID: user.ID, Name: req.Name, NoteSlug: req.NoteSlug,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toHabitResponse(*habit))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	targetUserID := r.URL.Query().Get("userId")

	habits, err := h.List.Execute(r.Context(), user.ID, targetUserID)
	if err == application.ErrNoCoachAccess {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]habitResponse, 0, len(habits))
	for _, hb := range habits {
		resp = append(resp, toHabitResponse(hb))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleUpdate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req updateHabitRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	habit, err := h.Update.Execute(r.Context(), application.UpdateHabitInput{
		UserID: user.ID, HabitID: id, Name: req.Name, NoteSlug: req.NoteSlug,
	})
	if err == application.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toHabitResponse(*habit))
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

func (h *Handler) HandleCheckIn(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req checkInRequest
	_ = httpkit.DecodeJSON(r, &req)

	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "date must be YYYY-MM-DD")
			return
		}
		date = parsed
	}

	err := h.CheckIn.Execute(r.Context(), user.ID, id, date)
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

func (h *Handler) HandleStats(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	stats, err := h.GetStats.Execute(r.Context(), user.ID, id)
	if err == application.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toStatsResponse(stats))
}

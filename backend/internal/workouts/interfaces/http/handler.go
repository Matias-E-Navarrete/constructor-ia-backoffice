package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
	"rimu/backend/internal/workouts/application"
	"rimu/backend/internal/workouts/domain"
)

type Handler struct {
	LogSession         *application.LogSession
	ListSessions       *application.ListSessions
	GetSession         *application.GetSession
	Delete             *application.DeleteSession
	GetProgress        *application.GetExerciseProgress
	GetPersonalRecords *application.GetPersonalRecords
	CreateRoutine      *application.CreateRoutine
	ListRoutines       *application.ListRoutines
	DeleteRoutine      *application.DeleteRoutine
}

func (h *Handler) HandleLogSession(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req logSessionRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	date, err := parseDate(req.SessionDate, time.Now())
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "session_date must be YYYY-MM-DD")
		return
	}

	sets := make([]domain.Set, 0, len(req.Sets))
	for _, s := range req.Sets {
		sets = append(sets, domain.Set{ExerciseName: s.ExerciseName, SetNumber: s.SetNumber, Reps: s.Reps, WeightKg: s.WeightKg})
	}

	session, err := h.LogSession.Execute(r.Context(), application.LogSessionInput{
		UserID: user.ID, SessionDate: date, Notes: req.Notes, NoteSlug: req.NoteSlug, Sets: sets,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toSessionResponse(*session))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	targetUserID := r.URL.Query().Get("userId")

	sessions, err := h.ListSessions.Execute(r.Context(), user.ID, targetUserID)
	if err == application.ErrNoCoachAccess {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]sessionResponse, 0, len(sessions))
	for _, s := range sessions {
		resp = append(resp, toSessionResponse(s))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleGet(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	session, err := h.GetSession.Execute(r.Context(), user.ID, id)
	if err == application.ErrNoCoachAccess {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toSessionResponse(*session))
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

func (h *Handler) HandleProgress(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	exercise := r.URL.Query().Get("exercise")
	if exercise == "" {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "missing_exercise", "exercise query param is required")
		return
	}

	progress, err := h.GetProgress.Execute(r.Context(), user.ID, exercise)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toProgressResponse(progress))
}

func (h *Handler) HandlePersonalRecords(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	prs, err := h.GetPersonalRecords.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]personalRecordResponse, 0, len(prs))
	for _, pr := range prs {
		resp = append(resp, toPersonalRecordResponse(pr))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleCreateRoutine(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createRoutineRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	routine, err := h.CreateRoutine.Execute(r.Context(), application.CreateRoutineInput{
		UserID: user.ID, Name: req.Name, Exercises: req.Exercises,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toRoutineResponse(*routine))
}

func (h *Handler) HandleListRoutines(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	routines, err := h.ListRoutines.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]routineResponse, 0, len(routines))
	for _, routine := range routines {
		resp = append(resp, toRoutineResponse(routine))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleDeleteRoutine(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.DeleteRoutine.Execute(r.Context(), user.ID, id); err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

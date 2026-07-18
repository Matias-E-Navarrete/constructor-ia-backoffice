package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
	"rimu/backend/internal/studies/application"
	"rimu/backend/internal/studies/domain"
)

type Handler struct {
	CreateSubject *application.CreateSubject
	ListSubjects  *application.ListSubjects
	DeleteSubject *application.DeleteSubject
	LogSession    *application.LogStudySession
	ListSessions  *application.ListSessions
	GetOverview   *application.GetStudyOverview
}

func (h *Handler) HandleCreateSubject(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createSubjectRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	subject, err := h.CreateSubject.Execute(r.Context(), application.CreateSubjectInput{
		UserID: user.ID, Name: req.Name, NoteSlug: req.NoteSlug,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toSubjectResponse(*subject))
}

func (h *Handler) HandleListSubjects(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	subjects, err := h.ListSubjects.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]subjectResponse, 0, len(subjects))
	for _, s := range subjects {
		resp = append(resp, toSubjectResponse(s))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleDeleteSubject(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.DeleteSubject.Execute(r.Context(), user.ID, id); err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
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

	session, err := h.LogSession.Execute(r.Context(), application.LogStudySessionInput{
		UserID: user.ID, SubjectID: req.SubjectID, SessionDate: date,
		DurationMinutes: req.DurationMinutes, Topic: req.Topic,
	})
	if err == domain.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toSessionResponse(*session))
}

func (h *Handler) HandleListSessions(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	subjectID := r.URL.Query().Get("subjectId")

	sessions, err := h.ListSessions.Execute(r.Context(), user.ID, subjectID)
	if err == domain.ErrNotOwner {
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

func (h *Handler) HandleOverview(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	overview, err := h.GetOverview.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]overviewResponse, 0, len(overview))
	for _, o := range overview {
		resp = append(resp, toOverviewResponse(o))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

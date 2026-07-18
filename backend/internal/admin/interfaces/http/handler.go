package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/admin/application"
	admindomain "rimu/backend/internal/admin/domain"
	"rimu/backend/internal/platform/httpkit"
	userdomain "rimu/backend/internal/user/domain"
)

type Handler struct {
	ListRoadmap         *application.ListRoadmap
	CreateRoadmapItem   *application.CreateRoadmapItem
	UpdateRoadmapStatus *application.UpdateRoadmapStatus

	ListFeatureFlags  *application.ListFeatureFlags
	ToggleFeatureFlag *application.ToggleFeatureFlag

	ListUsers   *application.ListUsers
	SetUserPlan *application.SetUserPlan
	SetUserRole *application.SetUserRole

	RunBackendTests  *application.RunBackendTests
	RunFrontendTests *application.RunFrontendTests
}

func (h *Handler) HandleListRoadmap(w nethttp.ResponseWriter, r *nethttp.Request) {
	items, err := h.ListRoadmap.Execute(r.Context())
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]roadmapItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, toRoadmapItemResponse(it))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleCreateRoadmapItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req createRoadmapItemRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	item, err := h.CreateRoadmapItem.Execute(r.Context(), application.CreateRoadmapItemInput{
		Title: req.Title, Description: req.Description, Kind: admindomain.Kind(req.Kind),
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toRoadmapItemResponse(*item))
}

func (h *Handler) HandleUpdateRoadmapStatus(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req updateRoadmapStatusRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	item, err := h.UpdateRoadmapStatus.Execute(r.Context(), id, admindomain.Status(req.Status))
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toRoadmapItemResponse(*item))
}

func (h *Handler) HandleListFeatureFlags(w nethttp.ResponseWriter, r *nethttp.Request) {
	flags, err := h.ListFeatureFlags.Execute(r.Context())
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]flagResponse, 0, len(flags))
	for _, f := range flags {
		resp = append(resp, toFlagResponse(f))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleToggleFeatureFlag(w nethttp.ResponseWriter, r *nethttp.Request) {
	key := chi.URLParam(r, "key")
	var req toggleFlagRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	flag, err := h.ToggleFeatureFlag.Execute(r.Context(), key, req.Enabled)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toFlagResponse(flag))
}

func (h *Handler) HandleListUsers(w nethttp.ResponseWriter, r *nethttp.Request) {
	users, err := h.ListUsers.Execute(r.Context())
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleSetUserPlan(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req setUserPlanRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	u, err := h.SetUserPlan.Execute(r.Context(), id, userdomain.Plan(req.Plan))
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toUserResponse(u))
}

func (h *Handler) HandleSetUserRole(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req setUserRoleRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	u, err := h.SetUserRole.Execute(r.Context(), id, userdomain.Role(req.Role))
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toUserResponse(u))
}

func (h *Handler) HandleRunTests(w nethttp.ResponseWriter, r *nethttp.Request) {
	suite := r.URL.Query().Get("suite")

	var result application.TestRunResult
	var err error
	switch suite {
	case "backend":
		result, err = h.RunBackendTests.Execute(r.Context())
	case "frontend":
		result, err = h.RunFrontendTests.Execute(r.Context())
	default:
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_suite", "suite must be 'backend' or 'frontend'")
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, result)
}

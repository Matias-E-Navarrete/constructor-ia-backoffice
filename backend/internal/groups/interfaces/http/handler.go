package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/groups/application"
	"rimu/backend/internal/groups/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	CreateFamily   *application.CreateFamilyGroup
	CreateCoaching *application.CreateCoachingGroup
	InviteMember   *application.InviteMember
	RemoveMember   *application.RemoveMember
	ListMyGroups   *application.ListMyGroups
	GetGroup       *application.GetGroup
}

func (h *Handler) HandleCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createGroupRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	switch domain.Kind(req.Kind) {
	case domain.KindFamily:
		g, err := h.CreateFamily.Execute(r.Context(), user.ID, req.Name)
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		httpkit.WriteJSON(w, nethttp.StatusCreated, toGroupResponse(*g))
	case domain.KindCoaching:
		g, err := h.CreateCoaching.Execute(r.Context(), user.ID, user.IsPro(), req.Name)
		if err == application.ErrProRequired {
			httpkit.WriteError(w, nethttp.StatusForbidden, "upgrade_required", err.Error())
			return
		}
		if err != nil {
			httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		httpkit.WriteJSON(w, nethttp.StatusCreated, toGroupResponse(*g))
	default:
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_kind", "kind must be 'family' or 'coaching'")
	}
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	memberships, err := h.ListMyGroups.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]membershipResponse, 0, len(memberships))
	for _, m := range memberships {
		resp = append(resp, membershipResponse{Group: toGroupResponse(m.Group), Role: string(m.Role)})
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleGet(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	groupID := chi.URLParam(r, "id")

	out, err := h.GetGroup.Execute(r.Context(), user.ID, groupID)
	if err == application.ErrNotMember {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}

	members := make([]memberResponse, 0, len(out.Members))
	for _, m := range out.Members {
		members = append(members, memberResponse{UserID: m.UserID, Role: string(m.Role)})
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, groupDetailResponse{Group: toGroupResponse(out.Group), Members: members})
}

func (h *Handler) HandleInvite(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	groupID := chi.URLParam(r, "id")

	var req inviteMemberRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	member, err := h.InviteMember.Execute(r.Context(), application.InviteMemberInput{
		RequesterUserID: user.ID,
		GroupID:         groupID,
		InviteeEmail:    req.Email,
	})
	if err == application.ErrNotAuthorized {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, memberResponse{UserID: member.UserID, Role: string(member.Role)})
}

func (h *Handler) HandleRemove(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	groupID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	err := h.RemoveMember.Execute(r.Context(), application.RemoveMemberInput{
		RequesterUserID: user.ID,
		GroupID:         groupID,
		TargetUserID:    targetUserID,
	})
	if err == application.ErrNotAuthorized {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

package http

import (
	nethttp "net/http"

	"rimu/backend/internal/assistant/application"
	"rimu/backend/internal/assistant/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Send  *application.SendMessage
	List  *application.ListMessages
	Clear *application.ClearConversation
}

func (h *Handler) HandleSend(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req sendMessageRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil || req.Content == "" {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "content is required")
		return
	}

	reply, err := h.Send.Execute(r.Context(), user.ID, req.Content)
	if err == domain.ErrNotConfigured {
		httpkit.WriteError(w, nethttp.StatusServiceUnavailable, "assistant_not_configured", "the assistant has no API key configured on this deployment")
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadGateway, "assistant_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toMessageResponse(*reply))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	messages, err := h.List.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]messageResponse, 0, len(messages))
	for _, m := range messages {
		resp = append(resp, toMessageResponse(m))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleClear(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	if err := h.Clear.Execute(r.Context(), user.ID); err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

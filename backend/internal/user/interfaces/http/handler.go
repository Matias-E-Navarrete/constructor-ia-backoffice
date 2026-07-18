package http

import (
	nethttp "net/http"

	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
	"rimu/backend/internal/user/application"
	"rimu/backend/internal/user/domain"
)

type Handler struct {
	Register *application.RegisterUser
	Login    *application.AuthenticateUser
	Upgrade  *application.UpgradePlan
}

func (h *Handler) HandleRegister(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req registerRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	u, err := h.Register.Execute(r.Context(), application.RegisterUserInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeUserError(w, err)
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toUserResponse(u))
}

func (h *Handler) HandleLogin(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req loginRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	out, err := h.Login.Execute(r.Context(), application.AuthenticateUserInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeUserError(w, err)
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, authResponse{Token: out.Token, User: toUserResponse(out.User)})
}

func (h *Handler) HandleMe(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	httpkit.WriteJSON(w, nethttp.StatusOK, toUserResponse(user))
}

func (h *Handler) HandleUpgrade(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	updated, err := h.Upgrade.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", "could not upgrade plan")
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toUserResponse(updated))
}

func writeUserError(w nethttp.ResponseWriter, err error) {
	switch err {
	case domain.ErrEmailTaken:
		httpkit.WriteError(w, nethttp.StatusConflict, "email_taken", err.Error())
	case application.ErrWeakPassword:
		httpkit.WriteError(w, nethttp.StatusBadRequest, "weak_password", err.Error())
	case application.ErrInvalidCredentials:
		httpkit.WriteError(w, nethttp.StatusUnauthorized, "invalid_credentials", err.Error())
	default:
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
	}
}

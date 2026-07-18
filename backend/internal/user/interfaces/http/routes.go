package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// Mount registers /api/auth/* (public) and /api/me* (authenticated) routes.
func Mount(r chi.Router, h *Handler, requireAuth func(nethttp.Handler) nethttp.Handler) {
	r.Post("/auth/register", h.HandleRegister)
	r.Post("/auth/login", h.HandleLogin)

	r.Group(func(r chi.Router) {
		r.Use(requireAuth)
		r.Get("/me", h.HandleMe)
		r.Post("/me/upgrade", h.HandleUpgrade)
		r.Patch("/me/planner-hours", h.HandlePlannerHours)
	})
}

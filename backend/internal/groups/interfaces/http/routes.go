package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// Mount registers /api/groups* behind RequireAuth + the "groups" feature flag.
func Mount(r chi.Router, h *Handler, requireAuth, requireFeature func(nethttp.Handler) nethttp.Handler) {
	r.Route("/groups", func(r chi.Router) {
		r.Use(requireAuth, requireFeature)
		r.Post("/", h.HandleCreate)
		r.Get("/", h.HandleList)
		r.Get("/{id}", h.HandleGet)
		r.Post("/{id}/members", h.HandleInvite)
		r.Delete("/{id}/members/{userId}", h.HandleRemove)
	})
}

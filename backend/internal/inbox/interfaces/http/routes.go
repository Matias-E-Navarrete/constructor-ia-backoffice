package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// Mount registers /api/inbox* behind auth + the "inbox" feature flag.
func Mount(r chi.Router, h *Handler, requireAuth, requireFeature func(nethttp.Handler) nethttp.Handler) {
	r.Route("/inbox", func(r chi.Router) {
		r.Use(requireAuth, requireFeature)

		r.Get("/", h.HandleList)
		r.Post("/", h.HandleCapture)
		r.Patch("/{id}/pin", h.HandleSetPinned)
		r.Delete("/{id}", h.HandleDelete)
	})
}

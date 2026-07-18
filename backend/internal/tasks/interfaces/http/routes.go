package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// Mount registers /api/tasks* behind auth + the "tasks" feature flag.
func Mount(r chi.Router, h *Handler, requireAuth, requireFeature func(nethttp.Handler) nethttp.Handler) {
	r.Route("/tasks", func(r chi.Router) {
		r.Use(requireAuth, requireFeature)

		r.Get("/", h.HandleList)
		r.Post("/", h.HandleCreate)
		r.Post("/reorder", h.HandleReorder)
		r.Patch("/{id}", h.HandleUpdate)
		r.Delete("/{id}", h.HandleDelete)
		r.Post("/{id}/complete", h.HandleComplete)
		r.Post("/{id}/uncomplete", h.HandleUncomplete)
		r.Patch("/{id}/status", h.HandleSetStatus)
		r.Patch("/{id}/quadrant", h.HandleSetQuadrant)
		r.Patch("/{id}/schedule", h.HandleSchedule)
	})
}

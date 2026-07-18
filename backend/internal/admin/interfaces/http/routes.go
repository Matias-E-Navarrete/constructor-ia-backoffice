package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// Mount registers /api/admin/* behind RequireAuth + RequireAdmin.
func Mount(r chi.Router, h *Handler, requireAuth, requireAdmin func(nethttp.Handler) nethttp.Handler) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(requireAuth, requireAdmin)

		r.Get("/roadmap", h.HandleListRoadmap)
		r.Post("/roadmap", h.HandleCreateRoadmapItem)
		r.Patch("/roadmap/{id}", h.HandleUpdateRoadmapStatus)

		r.Get("/feature-flags", h.HandleListFeatureFlags)
		r.Patch("/feature-flags/{key}", h.HandleToggleFeatureFlag)

		r.Get("/users", h.HandleListUsers)
		r.Patch("/users/{id}/plan", h.HandleSetUserPlan)
		r.Patch("/users/{id}/role", h.HandleSetUserRole)

		r.Post("/tests/run", h.HandleRunTests)
	})
}

package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth        func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag  func(nethttp.Handler) nethttp.Handler // "habits"
	RequireStatsFlag   func(nethttp.Handler) nethttp.Handler // "habits.stats"
	RequirePro         func(nethttp.Handler) nethttp.Handler
}

// Mount registers /api/habits* behind auth + the module kill switch; the
// stats endpoint additionally requires the pro plan and its own sub-feature flag.
func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/habits", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag)

		r.Get("/", h.HandleList)
		r.Post("/", h.HandleCreate)
		r.Patch("/{id}", h.HandleUpdate)
		r.Delete("/{id}", h.HandleDelete)
		r.Post("/{id}/log", h.HandleCheckIn)

		r.With(mw.RequireStatsFlag, mw.RequirePro).Get("/{id}/stats", h.HandleStats)
	})
}

package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth         func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag   func(nethttp.Handler) nethttp.Handler // "studies"
	RequireOverviewFlag func(nethttp.Handler) nethttp.Handler // "studies.overview"
	RequirePro          func(nethttp.Handler) nethttp.Handler
}

// Mount registers /api/studies* behind auth + the module kill switch; the
// overview endpoint additionally requires the pro plan and its own
// sub-feature flag.
func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/studies", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag)

		r.Get("/subjects", h.HandleListSubjects)
		r.Post("/subjects", h.HandleCreateSubject)
		r.Delete("/subjects/{id}", h.HandleDeleteSubject)

		r.Get("/sessions", h.HandleListSessions)
		r.Post("/sessions", h.HandleLogSession)

		r.With(mw.RequireOverviewFlag, mw.RequirePro).Get("/overview", h.HandleOverview)
	})
}

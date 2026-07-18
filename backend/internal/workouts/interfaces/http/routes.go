package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth          func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag    func(nethttp.Handler) nethttp.Handler // "workouts"
	RequireProgressFlag  func(nethttp.Handler) nethttp.Handler // "workouts.progress"
	RequireRoutinesFlag  func(nethttp.Handler) nethttp.Handler // "workouts.routines"
	RequirePro           func(nethttp.Handler) nethttp.Handler
}

func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/workouts", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag)

		r.Get("/", h.HandleList)
		r.Post("/", h.HandleLogSession)
		r.Get("/{id}", h.HandleGet)
		r.Delete("/{id}", h.HandleDelete)

		r.With(mw.RequireProgressFlag, mw.RequirePro).Get("/progress", h.HandleProgress)
		r.With(mw.RequireProgressFlag, mw.RequirePro).Get("/personal-records", h.HandlePersonalRecords)

		r.Route("/routines", func(r chi.Router) {
			r.Use(mw.RequireRoutinesFlag, mw.RequirePro)
			r.Get("/", h.HandleListRoutines)
			r.Post("/", h.HandleCreateRoutine)
			r.Delete("/{id}", h.HandleDeleteRoutine)
		})
	})
}

package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth        func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag  func(nethttp.Handler) nethttp.Handler // "notes"
	RequireGraphFlag   func(nethttp.Handler) nethttp.Handler // "notes.graph"
	RequireExportFlag  func(nethttp.Handler) nethttp.Handler // "notes.export"
	RequirePro         func(nethttp.Handler) nethttp.Handler
}

func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/notes", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag)

		r.Get("/", h.HandleList)
		r.Post("/", h.HandleCreate)
		r.Get("/{slug}", h.HandleGet)
		r.Patch("/{slug}", h.HandleUpdate)
		r.Delete("/{slug}", h.HandleDelete)

		r.With(mw.RequireGraphFlag, mw.RequirePro).Get("/graph", h.HandleGraph)
		r.With(mw.RequireExportFlag, mw.RequirePro).Get("/export", h.HandleExport)
	})
}

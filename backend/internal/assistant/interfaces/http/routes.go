package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth       func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag func(nethttp.Handler) nethttp.Handler // "assistant"
	RequirePro        func(nethttp.Handler) nethttp.Handler
}

// Mount registers /api/assistant* behind auth, the module kill switch, and
// the pro plan — the assistant is a premium capability end to end, not just
// one sub-feature within a free module.
func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/assistant", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag, mw.RequirePro)

		r.Get("/messages", h.HandleList)
		r.Post("/messages", h.HandleSend)
		r.Delete("/messages", h.HandleClear)
	})
}

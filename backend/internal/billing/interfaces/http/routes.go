package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth func(nethttp.Handler) nethttp.Handler
}

// Mount registers /api/billing*. The webhook route deliberately sits
// outside RequireAuth — Stripe calls it directly with its own signature,
// not a session belonging to any of our users.
func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/billing", func(r chi.Router) {
		r.Post("/webhook", h.HandleWebhook)

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAuth)
			r.Get("/config", h.HandleConfig)
			r.Post("/checkout", h.HandleCheckout)
		})
	})
}

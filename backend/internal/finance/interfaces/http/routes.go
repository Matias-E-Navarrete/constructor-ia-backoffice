package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Middlewares struct {
	RequireAuth        func(nethttp.Handler) nethttp.Handler
	RequireModuleFlag  func(nethttp.Handler) nethttp.Handler // "finance"
	RequireSummaryFlag func(nethttp.Handler) nethttp.Handler // "finance.summary"
	RequireExportFlag  func(nethttp.Handler) nethttp.Handler // "finance.export"
	RequirePro         func(nethttp.Handler) nethttp.Handler
}

func Mount(r chi.Router, h *Handler, mw Middlewares) {
	r.Route("/finance", func(r chi.Router) {
		r.Use(mw.RequireAuth, mw.RequireModuleFlag)

		r.Get("/transactions", h.HandleList)
		r.Post("/transactions", h.HandleRecord)
		r.Delete("/transactions/{id}", h.HandleDelete)

		r.Get("/accounts", h.HandleListAccounts)
		r.Post("/accounts", h.HandleCreateAccount)
		r.Get("/convert", h.HandleConvert)

		r.With(mw.RequireSummaryFlag, mw.RequirePro).Get("/summary", h.HandleSummary)
		r.With(mw.RequireSummaryFlag, mw.RequirePro).Get("/upcoming", h.HandleUpcoming)
		r.With(mw.RequireExportFlag, mw.RequirePro).Get("/export", h.HandleExport)
		r.With(mw.RequireExportFlag, mw.RequirePro).Post("/export/pdf", h.HandleExportPDF)
	})
}

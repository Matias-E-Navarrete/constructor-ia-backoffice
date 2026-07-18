package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/finance/application"
	"rimu/backend/internal/finance/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Record       *application.RecordTransaction
	List         *application.ListTransactions
	Delete       *application.DeleteTransaction
	GetSummary   *application.GetFinanceSummary
	ExportCSV    *application.ExportTransactions
}

func (h *Handler) HandleRecord(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req recordTransactionRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	txDate, err := time.Parse("2006-01-02", req.TxDate)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_date", "tx_date must be YYYY-MM-DD")
		return
	}

	t, err := h.Record.Execute(r.Context(), application.RecordTransactionInput{
		UserID:      user.ID,
		Type:        domain.TxType(req.Type),
		AmountValue: req.Amount,
		Category:    req.Category,
		Description: req.Description,
		TxDate:      txDate,
		NoteSlug:    req.NoteSlug,
		GroupID:     req.GroupID,
	})
	if err == application.ErrNotGroupMember {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toTransactionResponse(*t))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var groupID *string
	if g := r.URL.Query().Get("groupId"); g != "" {
		groupID = &g
	}

	txs, err := h.List.Execute(r.Context(), user.ID, groupID)
	if err == application.ErrNotGroupMember {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	resp := make([]transactionResponse, 0, len(txs))
	for _, t := range txs {
		resp = append(resp, toTransactionResponse(t))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleDelete(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	err := h.Delete.Execute(r.Context(), user.ID, id)
	if err == application.ErrNotOwner {
		httpkit.WriteError(w, nethttp.StatusForbidden, "forbidden", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *Handler) HandleSummary(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	summary, err := h.GetSummary.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toSummaryResponse(summary))
}

func (h *Handler) HandleExport(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	csvBytes, err := h.ExportCSV.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")
	w.WriteHeader(nethttp.StatusOK)
	w.Write(csvBytes)
}

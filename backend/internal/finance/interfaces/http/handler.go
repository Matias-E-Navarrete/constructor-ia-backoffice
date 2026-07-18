package http

import (
	nethttp "net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/finance/application"
	"rimu/backend/internal/finance/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Record          *application.RecordTransaction
	List            *application.ListTransactions
	Delete          *application.DeleteTransaction
	GetSummary      *application.GetFinanceSummary
	ExportCSV       *application.ExportTransactions
	CreateAccount   *application.CreateAccount
	ListAccounts    *application.ListAccounts
	Convert         *application.ConvertCurrency
	GetUpcoming     *application.GetUpcomingBills
	ExportPDF       *application.ExportPDF
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
		UserID:             user.ID,
		Type:               domain.TxType(req.Type),
		AmountValue:        req.Amount,
		Category:           req.Category,
		Description:        req.Description,
		TxDate:             txDate,
		NoteSlug:           req.NoteSlug,
		GroupID:            req.GroupID,
		AccountID:          req.AccountID,
		Currency:           req.Currency,
		ExchangeRate:       req.ExchangeRate,
		Method:             req.Method,
		InstallmentsTotal:  req.InstallmentsTotal,
		InstallmentNumber:  req.InstallmentNumber,
		Recurring:          req.Recurring,
		RecurrenceInterval: req.RecurrenceInterval,
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

func (h *Handler) HandleCreateAccount(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createAccountRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}
	a, err := h.CreateAccount.Execute(r.Context(), application.CreateAccountInput{
		UserID: user.ID, Name: req.Name, Type: req.Type, Currency: req.Currency, InitialBalance: req.InitialBalance,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toAccountResponse(*a))
}

func (h *Handler) HandleListAccounts(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	accounts, err := h.ListAccounts.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]accountResponse, 0, len(accounts))
	for _, a := range accounts {
		resp = append(resp, toAccountResponse(a))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleConvert(w nethttp.ResponseWriter, r *nethttp.Request) {
	q := r.URL.Query()
	amount := parseFloat(q.Get("amount"))
	from := q.Get("from")
	to := q.Get("to")

	result, err := h.Convert.Execute(r.Context(), amount, from, to)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "unknown_currency", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, convertResponse{ConvertedAmount: result.ConvertedAmount, Rate: result.Rate})
}

func (h *Handler) HandleUpcoming(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	q := r.URL.Query()

	from, to := monthRange(time.Now())
	if f := q.Get("from"); f != "" {
		if parsed, err := time.Parse("2006-01-02", f); err == nil {
			from = parsed
		}
	}
	if t := q.Get("to"); t != "" {
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			to = parsed
		}
	}

	bills, err := h.GetUpcoming.Execute(r.Context(), user.ID, from, to)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]upcomingBillResponse, 0, len(bills))
	for _, b := range bills {
		resp = append(resp, upcomingBillResponse{
			Transaction: toTransactionResponse(b.Transaction),
			NextDate:    b.NextDate.Format("2006-01-02"),
		})
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleExportPDF(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req exportPDFRequest
	_ = httpkit.DecodeJSON(r, &req)

	in := application.ExportPDFInput{UserID: user.ID, Type: req.Type, Categories: req.Categories, AccountIDs: req.AccountIDs}
	if req.From != "" {
		if parsed, err := time.Parse("2006-01-02", req.From); err == nil {
			in.From = parsed
		}
	}
	if req.To != "" {
		if parsed, err := time.Parse("2006-01-02", req.To); err == nil {
			in.To = parsed
		}
	}

	pdfBytes, err := h.ExportPDF.Execute(r.Context(), in)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=informe.pdf")
	w.WriteHeader(nethttp.StatusOK)
	w.Write(pdfBytes)
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func monthRange(now time.Time) (time.Time, time.Time) {
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 1, -1)
	return from, to
}

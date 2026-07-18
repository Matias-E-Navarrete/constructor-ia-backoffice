package http

import (
	"rimu/backend/internal/finance/domain"
)

type recordTransactionRequest struct {
	Type               string   `json:"type"`
	Amount             float64  `json:"amount"`
	Category           string   `json:"category"`
	Description        string   `json:"description,omitempty"`
	TxDate             string   `json:"tx_date"` // YYYY-MM-DD
	NoteSlug           *string  `json:"note_slug,omitempty"`
	GroupID            *string  `json:"group_id,omitempty"`
	AccountID          *string  `json:"account_id,omitempty"`
	Currency           string   `json:"currency,omitempty"`
	ExchangeRate       *float64 `json:"exchange_rate,omitempty"`
	Method             string   `json:"method,omitempty"`
	InstallmentsTotal  *int     `json:"installments_total,omitempty"`
	InstallmentNumber  *int     `json:"installment_number,omitempty"`
	Recurring          bool     `json:"recurring,omitempty"`
	RecurrenceInterval *string  `json:"recurrence_interval,omitempty"`
}

type transactionResponse struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	Amount             float64  `json:"amount"`
	Category           string   `json:"category"`
	Description        string   `json:"description"`
	TxDate             string   `json:"tx_date"`
	NoteSlug           *string  `json:"note_slug,omitempty"`
	GroupID            *string  `json:"group_id,omitempty"`
	AccountID          *string  `json:"account_id,omitempty"`
	Currency           string   `json:"currency"`
	ExchangeRate       *float64 `json:"exchange_rate,omitempty"`
	Method             string   `json:"method"`
	InstallmentsTotal  *int     `json:"installments_total,omitempty"`
	InstallmentNumber  *int     `json:"installment_number,omitempty"`
	Recurring          bool     `json:"recurring"`
	RecurrenceInterval *string  `json:"recurrence_interval,omitempty"`
}

func toTransactionResponse(t domain.Transaction) transactionResponse {
	return transactionResponse{
		ID:                 t.ID,
		Type:               string(t.Type),
		Amount:             t.Amount.Value(),
		Category:           t.Category,
		Description:        t.Description,
		TxDate:             t.TxDate.Format("2006-01-02"),
		NoteSlug:           t.NoteSlug,
		GroupID:            t.GroupID,
		AccountID:          t.AccountID,
		Currency:           t.Currency,
		ExchangeRate:       t.ExchangeRate,
		Method:             t.Method,
		InstallmentsTotal:  t.InstallmentsTotal,
		InstallmentNumber:  t.InstallmentNumber,
		Recurring:          t.Recurring,
		RecurrenceInterval: t.RecurrenceInterval,
	}
}

type categoryBreakdownResponse struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type summaryResponse struct {
	Balance    float64                     `json:"balance"`
	Income     float64                     `json:"income"`
	Expenses   float64                     `json:"expenses"`
	ByCategory []categoryBreakdownResponse `json:"by_category"`
}

func toSummaryResponse(s domain.Summary) summaryResponse {
	byCategory := make([]categoryBreakdownResponse, 0, len(s.ByCategory))
	for _, c := range s.ByCategory {
		byCategory = append(byCategory, categoryBreakdownResponse{Category: c.Category, Total: c.Total})
	}
	return summaryResponse{Balance: s.Balance, Income: s.Income, Expenses: s.Expenses, ByCategory: byCategory}
}

type createAccountRequest struct {
	Name           string  `json:"name"`
	Type           string  `json:"type,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	InitialBalance float64 `json:"initial_balance,omitempty"`
}

type accountResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Currency       string  `json:"currency"`
	InitialBalance float64 `json:"initial_balance"`
	Active         bool    `json:"active"`
}

func toAccountResponse(a domain.Account) accountResponse {
	return accountResponse{
		ID: a.ID, Name: a.Name, Type: a.Type, Currency: a.Currency,
		InitialBalance: a.InitialBalance, Active: a.Active,
	}
}

type convertResponse struct {
	ConvertedAmount float64 `json:"converted_amount"`
	Rate            float64 `json:"rate"`
}

type upcomingBillResponse struct {
	Transaction transactionResponse `json:"transaction"`
	NextDate    string              `json:"next_date"`
}

type exportPDFRequest struct {
	Type       string   `json:"type,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	Categories []string `json:"categories,omitempty"`
	AccountIDs []string `json:"account_ids,omitempty"`
}

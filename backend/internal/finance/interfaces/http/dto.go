package http

import (
	"rimu/backend/internal/finance/domain"
)

type recordTransactionRequest struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
	TxDate      string  `json:"tx_date"` // YYYY-MM-DD
	NoteSlug    *string `json:"note_slug,omitempty"`
	GroupID     *string `json:"group_id,omitempty"`
}

type transactionResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	TxDate      string  `json:"tx_date"`
	NoteSlug    *string `json:"note_slug,omitempty"`
	GroupID     *string `json:"group_id,omitempty"`
}

func toTransactionResponse(t domain.Transaction) transactionResponse {
	return transactionResponse{
		ID:          t.ID,
		Type:        string(t.Type),
		Amount:      t.Amount.Value(),
		Category:    t.Category,
		Description: t.Description,
		TxDate:      t.TxDate.Format("2006-01-02"),
		NoteSlug:    t.NoteSlug,
		GroupID:     t.GroupID,
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

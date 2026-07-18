package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestFinance_RecordListSummaryExportGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	tx := map[string]interface{}{
		"type": "expense", "amount": 50, "category": "comida", "tx_date": "2026-07-01",
	}
	var created struct{ ID string }
	status := s.Do(http.MethodPost, "/api/finance/transactions", token, tx, &created)
	if status != http.StatusCreated || created.ID == "" {
		t.Fatalf("record transaction: status=%d body=%+v", status, created)
	}

	var list []struct{ ID string }
	status = s.Do(http.MethodGet, "/api/finance/transactions", token, nil, &list)
	if status != http.StatusOK || len(list) != 1 {
		t.Fatalf("list transactions: status=%d len=%d", status, len(list))
	}

	status = s.Do(http.MethodGet, "/api/finance/summary", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("summary on free plan: expected 403, got %d", status)
	}
	status = s.Do(http.MethodGet, "/api/finance/export", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("export on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var summary struct {
		Balance float64 `json:"balance"`
	}
	status = s.Do(http.MethodGet, "/api/finance/summary", token, nil, &summary)
	if status != http.StatusOK || summary.Balance != -50 {
		t.Fatalf("summary on pro plan: status=%d body=%+v", status, summary)
	}

	status = s.Do(http.MethodGet, "/api/finance/export", token, nil, nil)
	if status != http.StatusOK {
		t.Fatalf("export on pro plan: expected 200, got %d", status)
	}
}

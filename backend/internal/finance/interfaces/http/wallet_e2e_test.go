package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestWallet_AccountsCurrencyInstallmentsRecurringExportGating(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var account struct {
		ID       string
		Currency string
	}
	status := s.Do(http.MethodPost, "/api/finance/accounts", token, map[string]interface{}{
		"name": "Bank of America Checking", "currency": "USD", "initial_balance": 1000,
	}, &account)
	if status != http.StatusCreated || account.ID == "" {
		t.Fatalf("create account: status=%d body=%+v", status, account)
	}

	var accounts []struct{ ID string }
	s.Do(http.MethodGet, "/api/finance/accounts", token, nil, &accounts)
	if len(accounts) != 1 {
		t.Fatalf("list accounts: expected 1, got %d", len(accounts))
	}

	var conv struct {
		ConvertedAmount float64 `json:"converted_amount"`
		Rate            float64 `json:"rate"`
	}
	status = s.Do(http.MethodGet, "/api/finance/convert?amount=20&from=GBP&to=USD", token, nil, &conv)
	if status != http.StatusOK || conv.ConvertedAmount <= 0 {
		t.Fatalf("convert: status=%d body=%+v", status, conv)
	}

	var tx struct {
		ID                string
		InstallmentsTotal int `json:"installments_total"`
		Recurring         bool
	}
	status = s.Do(http.MethodPost, "/api/finance/transactions", token, map[string]interface{}{
		"type": "expense", "amount": 212.66, "category": "Pago de factura", "tx_date": "2026-04-22",
		"account_id": account.ID, "currency": "USD", "method": "Amex Platinum",
		"installments_total": 3, "installment_number": 1, "recurring": true, "recurrence_interval": "monthly",
	}, &tx)
	if status != http.StatusCreated || !tx.Recurring || tx.InstallmentsTotal != 3 {
		t.Fatalf("record transaction: status=%d body=%+v", status, tx)
	}

	status = s.Do(http.MethodGet, "/api/finance/upcoming?from=2026-05-01&to=2026-05-31", token, nil, nil)
	if status != http.StatusForbidden {
		t.Fatalf("upcoming on free plan: expected 403, got %d", status)
	}
	status = s.Do(http.MethodPost, "/api/finance/export/pdf", token, map[string]string{"type": "all"}, nil)
	if status != http.StatusForbidden {
		t.Fatalf("export pdf on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	var upcoming []struct {
		NextDate string `json:"next_date"`
	}
	status = s.Do(http.MethodGet, "/api/finance/upcoming?from=2026-05-01&to=2026-05-31", token, nil, &upcoming)
	if status != http.StatusOK || len(upcoming) != 1 || upcoming[0].NextDate != "2026-05-22" {
		t.Fatalf("upcoming on pro plan: status=%d body=%+v", status, upcoming)
	}

	status = s.Do(http.MethodPost, "/api/finance/export/pdf", token, map[string]string{"type": "all"}, nil)
	if status != http.StatusOK {
		t.Fatalf("export pdf on pro plan: expected 200, got %d", status)
	}
}

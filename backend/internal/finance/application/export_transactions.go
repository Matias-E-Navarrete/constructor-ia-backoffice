package application

import (
	"bytes"
	"context"
	"encoding/csv"
	"strconv"

	"rimu/backend/internal/finance/domain"
)

// ExportTransactions is a [PRO] use-case: renders every transaction the user
// owns as CSV bytes, ready to download.
type ExportTransactions struct {
	Repo domain.Repository
}

func (uc *ExportTransactions) Execute(ctx context.Context, userID string) ([]byte, error) {
	txs, err := uc.Repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Write([]string{"date", "type", "amount", "category", "description"})
	for _, t := range txs {
		w.Write([]string{
			t.TxDate.Format("2006-01-02"),
			string(t.Type),
			strconv.FormatFloat(t.Amount.Value(), 'f', 2, 64),
			t.Category,
			t.Description,
		})
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

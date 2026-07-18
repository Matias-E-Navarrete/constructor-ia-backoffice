package application

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"

	"rimu/backend/internal/finance/domain"
)

// ExportPDFInput mirrors the "Exportar informe" filters: transaction type,
// date range, categories and accounts (empty slices mean "all").
type ExportPDFInput struct {
	UserID     string
	Type       string // "all" | "income" | "expense"
	From       time.Time
	To         time.Time
	Categories []string
	AccountIDs []string
}

// ExportPDF is a [PRO] use-case: renders a filtered transaction report as a
// simple PDF (title, filter summary, table, totals).
type ExportPDF struct {
	Repo domain.Repository
}

func (uc *ExportPDF) Execute(ctx context.Context, in ExportPDFInput) ([]byte, error) {
	txs, err := uc.Repo.ListByUser(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	categorySet := toSet(in.Categories)
	accountSet := toSet(in.AccountIDs)

	filtered := make([]domain.Transaction, 0, len(txs))
	for _, t := range txs {
		if in.Type == "income" && t.Type != domain.TypeIncome {
			continue
		}
		if in.Type == "expense" && t.Type != domain.TypeExpense {
			continue
		}
		if !in.From.IsZero() && t.TxDate.Before(in.From) {
			continue
		}
		if !in.To.IsZero() && t.TxDate.After(in.To) {
			continue
		}
		if len(categorySet) > 0 && !categorySet[t.Category] {
			continue
		}
		if len(accountSet) > 0 {
			if t.AccountID == nil || !accountSet[*t.AccountID] {
				continue
			}
		}
		filtered = append(filtered, t)
	}

	return renderPDF(filtered), nil
}

func toSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	return set
}

func renderPDF(txs []domain.Transaction) []byte {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Informe financiero")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 8, "Fecha", "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, "Tipo", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, "Categoria", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 8, "Descripcion", "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, "Monto", "1", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	var total float64
	for _, t := range txs {
		pdf.CellFormat(30, 8, t.TxDate.Format("2006-01-02"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, string(t.Type), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 8, t.Category, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, t.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, formatAmount(t.Amount.Value()), "1", 1, "R", false, 0, "")
		total += t.SignedAmount()
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Total: "+formatAmount(total))

	var buf bytes.Buffer
	_ = pdf.Output(&buf)
	return buf.Bytes()
}

func formatAmount(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

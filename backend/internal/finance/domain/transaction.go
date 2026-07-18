package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("finance: not found")
var ErrInvalidAmount = errors.New("finance: amount must be greater than zero")

type TxType string

const (
	TypeIncome  TxType = "income"
	TypeExpense TxType = "expense"
)

// Amount is a value object: a validated, always-positive money quantity.
// Whether it adds to or subtracts from the balance depends on the
// transaction's Type (see Transaction.SignedAmount).
type Amount struct {
	value float64
}

func NewAmount(value float64) (Amount, error) {
	if value <= 0 {
		return Amount{}, ErrInvalidAmount
	}
	return Amount{value: value}, nil
}

func (a Amount) Value() float64 { return a.value }

type Transaction struct {
	ID                 string
	UserID             string
	Type               TxType
	Amount             Amount
	Category           string
	Description        string
	TxDate             time.Time
	NoteSlug           *string
	GroupID            *string
	AccountID          *string
	Currency           string
	ExchangeRate       *float64
	Method             string
	InstallmentsTotal  *int
	InstallmentNumber  *int
	Recurring          bool
	RecurrenceInterval *string
	CreatedAt          time.Time
}

// SignedAmount is positive for income, negative for expense — the quantity
// that should be summed to compute a balance.
func (t Transaction) SignedAmount() float64 {
	if t.Type == TypeExpense {
		return -t.Amount.Value()
	}
	return t.Amount.Value()
}

type CategoryBreakdown struct {
	Category string
	Total    float64
}

type Summary struct {
	Balance    float64
	Income     float64
	Expenses   float64
	ByCategory []CategoryBreakdown
}

// ComputeSummary is a pure function so the [PRO] summary use-case is
// trivially unit-testable without a database, mirroring habits.ComputeStats.
func ComputeSummary(transactions []Transaction) Summary {
	s := Summary{}
	byCategory := make(map[string]float64)

	for _, t := range transactions {
		switch t.Type {
		case TypeIncome:
			s.Income += t.Amount.Value()
		case TypeExpense:
			s.Expenses += t.Amount.Value()
		}
		s.Balance += t.SignedAmount()
		byCategory[t.Category] += t.SignedAmount()
	}

	for category, total := range byCategory {
		s.ByCategory = append(s.ByCategory, CategoryBreakdown{Category: category, Total: total})
	}
	return s
}

// UpcomingBill is a recurring transaction projected forward to its next
// occurrence date, for the "próximas cuentas" view.
type UpcomingBill struct {
	Transaction Transaction
	NextDate    time.Time
}

// ProjectUpcoming advances every recurring transaction by its interval
// until the projected date falls within [from, to], for calendar-style
// "upcoming bills" views. A pure function, mirroring ComputeSummary.
func ProjectUpcoming(transactions []Transaction, from, to time.Time) []UpcomingBill {
	var bills []UpcomingBill
	for _, t := range transactions {
		if !t.Recurring || t.RecurrenceInterval == nil {
			continue
		}
		next := t.TxDate
		for i := 0; i < 240 && next.Before(from); i++ { // cap iterations, ~20 years monthly
			next = advance(next, *t.RecurrenceInterval)
		}
		if !next.After(to) && !next.Before(from) {
			bills = append(bills, UpcomingBill{Transaction: t, NextDate: next})
		}
	}
	return bills
}

func advance(t time.Time, interval string) time.Time {
	switch interval {
	case "weekly":
		return t.AddDate(0, 0, 7)
	default: // monthly
		return t.AddDate(0, 1, 0)
	}
}

type Repository interface {
	Create(ctx context.Context, t *Transaction) error
	FindByID(ctx context.Context, id string) (*Transaction, error)
	ListByUser(ctx context.Context, userID string) ([]Transaction, error)
	ListByGroup(ctx context.Context, groupID string) ([]Transaction, error)
	Delete(ctx context.Context, id string) error
}

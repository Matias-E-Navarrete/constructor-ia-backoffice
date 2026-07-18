package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/finance/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const txColumns = `id, user_id, type, amount, category, description, tx_date, note_slug, group_id,
	account_id, currency, exchange_rate, method, installments_total, installment_number,
	recurring, recurrence_interval, created_at`

func (r *PostgresRepository) Create(ctx context.Context, t *domain.Transaction) error {
	if t.Currency == "" {
		t.Currency = "USD"
	}
	return r.pool.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, category, description, tx_date, note_slug, group_id,
			account_id, currency, exchange_rate, method, installments_total, installment_number,
			recurring, recurrence_interval)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at
	`, t.UserID, t.Type, t.Amount.Value(), t.Category, t.Description, t.TxDate, t.NoteSlug, t.GroupID,
		t.AccountID, t.Currency, t.ExchangeRate, t.Method, t.InstallmentsTotal, t.InstallmentNumber,
		t.Recurring, t.RecurrenceInterval,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return r.scanOne(r.pool.QueryRow(ctx, `SELECT `+txColumns+` FROM transactions WHERE id = $1`, id))
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+txColumns+` FROM transactions WHERE user_id = $1 ORDER BY tx_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return r.scanMany(rows)
}

func (r *PostgresRepository) ListByGroup(ctx context.Context, groupID string) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+txColumns+` FROM transactions WHERE group_id = $1 ORDER BY tx_date DESC
	`, groupID)
	if err != nil {
		return nil, err
	}
	return r.scanMany(rows)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) scanOne(row pgx.Row) (*domain.Transaction, error) {
	t := &domain.Transaction{}
	var amount float64
	err := row.Scan(&t.ID, &t.UserID, &t.Type, &amount, &t.Category, &t.Description, &t.TxDate, &t.NoteSlug, &t.GroupID,
		&t.AccountID, &t.Currency, &t.ExchangeRate, &t.Method, &t.InstallmentsTotal, &t.InstallmentNumber,
		&t.Recurring, &t.RecurrenceInterval, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	t.Amount, _ = domain.NewAmount(amount)
	return t, nil
}

func (r *PostgresRepository) scanMany(rows pgx.Rows) ([]domain.Transaction, error) {
	defer rows.Close()

	var txs []domain.Transaction
	for rows.Next() {
		t := domain.Transaction{}
		var amount float64
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &amount, &t.Category, &t.Description, &t.TxDate, &t.NoteSlug, &t.GroupID,
			&t.AccountID, &t.Currency, &t.ExchangeRate, &t.Method, &t.InstallmentsTotal, &t.InstallmentNumber,
			&t.Recurring, &t.RecurrenceInterval, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Amount, _ = domain.NewAmount(amount)
		txs = append(txs, t)
	}
	return txs, rows.Err()
}

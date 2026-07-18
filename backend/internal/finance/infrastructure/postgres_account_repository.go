package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/finance/domain"
)

type PostgresAccountRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAccountRepository(pool *pgxpool.Pool) *PostgresAccountRepository {
	return &PostgresAccountRepository{pool: pool}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, a *domain.Account) error {
	if a.Type == "" {
		a.Type = "Cuenta"
	}
	if a.Currency == "" {
		a.Currency = "USD"
	}
	return r.pool.QueryRow(ctx, `
		INSERT INTO accounts (user_id, name, type, currency, initial_balance, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, a.UserID, a.Name, a.Type, a.Currency, a.InitialBalance, a.Active).Scan(&a.ID, &a.CreatedAt)
}

func (r *PostgresAccountRepository) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	a := &domain.Account{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, type, currency, initial_balance, active, created_at
		FROM accounts WHERE id = $1
	`, id).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.InitialBalance, &a.Active, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	return a, err
}

func (r *PostgresAccountRepository) ListByUser(ctx context.Context, userID string) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, type, currency, initial_balance, active, created_at
		FROM accounts WHERE user_id = $1 ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		a := domain.Account{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Currency, &a.InitialBalance, &a.Active, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *PostgresAccountRepository) SetActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE accounts SET active = $2 WHERE id = $1`, id, active)
	return err
}

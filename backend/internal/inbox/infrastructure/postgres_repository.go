package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/inbox/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, item *domain.Item) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO inbox_items (user_id, content, pinned)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, item.UserID, item.Content, item.Pinned).Scan(&item.ID, &item.CreatedAt)
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, content, pinned, created_at FROM inbox_items
		WHERE user_id = $1 ORDER BY pinned DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		it := domain.Item{}
		if err := rows.Scan(&it.ID, &it.UserID, &it.Content, &it.Pinned, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Item, error) {
	it := &domain.Item{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, content, pinned, created_at FROM inbox_items WHERE id = $1
	`, id).Scan(&it.ID, &it.UserID, &it.Content, &it.Pinned, &it.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return it, err
}

func (r *PostgresRepository) SetPinned(ctx context.Context, id string, pinned bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE inbox_items SET pinned = $2 WHERE id = $1`, id, pinned)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inbox_items WHERE id = $1`, id)
	return err
}

package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/platform/featureflags/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context) ([]domain.Flag, error) {
	rows, err := r.pool.Query(ctx, `SELECT key, enabled, description FROM feature_flags ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []domain.Flag
	for rows.Next() {
		var f domain.Flag
		if err := rows.Scan(&f.Key, &f.Enabled, &f.Description); err != nil {
			return nil, err
		}
		flags = append(flags, f)
	}
	return flags, rows.Err()
}

func (r *PostgresRepository) SetEnabled(ctx context.Context, key string, enabled bool) (domain.Flag, error) {
	var f domain.Flag
	err := r.pool.QueryRow(ctx, `
		UPDATE feature_flags SET enabled = $2, updated_at = now()
		WHERE key = $1
		RETURNING key, enabled, description
	`, key, enabled).Scan(&f.Key, &f.Enabled, &f.Description)
	return f, err
}

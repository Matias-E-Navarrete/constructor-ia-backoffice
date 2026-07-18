package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/admin/domain"
)

type PostgresRoadmapRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRoadmapRepository(pool *pgxpool.Pool) *PostgresRoadmapRepository {
	return &PostgresRoadmapRepository{pool: pool}
}

func (r *PostgresRoadmapRepository) Create(ctx context.Context, item *domain.RoadmapItem) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO roadmap_items (title, description, kind, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Title, item.Description, item.Kind, item.Status).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *PostgresRoadmapRepository) List(ctx context.Context) ([]domain.RoadmapItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, kind, status, created_at, updated_at
		FROM roadmap_items ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.RoadmapItem
	for rows.Next() {
		it := domain.RoadmapItem{}
		if err := rows.Scan(&it.ID, &it.Title, &it.Description, &it.Kind, &it.Status, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *PostgresRoadmapRepository) UpdateStatus(ctx context.Context, id string, status domain.Status) (*domain.RoadmapItem, error) {
	it := &domain.RoadmapItem{}
	err := r.pool.QueryRow(ctx, `
		UPDATE roadmap_items SET status = $2, updated_at = now() WHERE id = $1
		RETURNING id, title, description, kind, status, created_at, updated_at
	`, id, status).Scan(&it.ID, &it.Title, &it.Description, &it.Kind, &it.Status, &it.CreatedAt, &it.UpdatedAt)
	return it, err
}

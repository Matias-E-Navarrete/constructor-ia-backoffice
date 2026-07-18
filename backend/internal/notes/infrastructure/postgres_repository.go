package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/notes/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, n *domain.Note) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO notes (user_id, slug, title, body_markdown, kind)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, n.UserID, n.Slug, n.Title, n.BodyMarkdown, n.Kind).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrSlugTaken
	}
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, n *domain.Note) error {
	return r.pool.QueryRow(ctx, `
		UPDATE notes SET title = $3, body_markdown = $4, updated_at = now()
		WHERE user_id = $1 AND slug = $2
		RETURNING updated_at
	`, n.UserID, n.Slug, n.Title, n.BodyMarkdown).Scan(&n.UpdatedAt)
}

func (r *PostgresRepository) FindBySlug(ctx context.Context, userID, slug string) (*domain.Note, error) {
	n := &domain.Note{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, slug, title, body_markdown, kind, created_at, updated_at
		FROM notes WHERE user_id = $1 AND slug = $2
	`, userID, slug).Scan(&n.ID, &n.UserID, &n.Slug, &n.Title, &n.BodyMarkdown, &n.Kind, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return n, err
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Note, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, slug, title, body_markdown, kind, created_at, updated_at
		FROM notes WHERE user_id = $1 ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		n := domain.Note{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Slug, &n.Title, &n.BodyMarkdown, &n.Kind, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *PostgresRepository) Delete(ctx context.Context, userID, slug string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notes WHERE user_id = $1 AND slug = $2`, userID, slug)
	return err
}

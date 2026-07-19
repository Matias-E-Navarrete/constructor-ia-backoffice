package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/assistant/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) SaveMessage(ctx context.Context, m *domain.Message) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO assistant_messages (user_id, role, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, m.UserID, string(m.Role), m.Content).Scan(&m.ID, &m.CreatedAt)
}

func (r *PostgresRepository) ListMessages(ctx context.Context, userID string) ([]domain.Message, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, role, content, created_at
		FROM assistant_messages WHERE user_id = $1 ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		var role string
		if err := rows.Scan(&m.ID, &m.UserID, &role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.Role = domain.Role(role)
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *PostgresRepository) ClearMessages(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM assistant_messages WHERE user_id = $1`, userID)
	return err
}

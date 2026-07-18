package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/habits/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, h *domain.Habit) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO habits (user_id, name, note_slug)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, h.UserID, h.Name, h.NoteSlug).Scan(&h.ID, &h.CreatedAt)
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Habit, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, note_slug, created_at, archived_at
		FROM habits WHERE user_id = $1 ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []domain.Habit
	for rows.Next() {
		h := domain.Habit{}
		if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &h.NoteSlug, &h.CreatedAt, &h.ArchivedAt); err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}
	return habits, rows.Err()
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Habit, error) {
	h := &domain.Habit{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, note_slug, created_at, archived_at
		FROM habits WHERE id = $1
	`, id).Scan(&h.ID, &h.UserID, &h.Name, &h.NoteSlug, &h.CreatedAt, &h.ArchivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return h, err
}

func (r *PostgresRepository) Update(ctx context.Context, h *domain.Habit) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE habits SET name = $2, note_slug = $3 WHERE id = $1
	`, h.ID, h.Name, h.NoteSlug)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM habits WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) UpsertLog(ctx context.Context, log domain.HabitLog) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO habit_logs (habit_id, log_date, completed)
		VALUES ($1, $2, $3)
		ON CONFLICT (habit_id, log_date) DO UPDATE SET completed = EXCLUDED.completed
	`, log.HabitID, log.LogDate, log.Completed)
	return err
}

func (r *PostgresRepository) ListLogs(ctx context.Context, habitID string) ([]domain.HabitLog, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT habit_id, log_date, completed FROM habit_logs WHERE habit_id = $1 ORDER BY log_date
	`, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.HabitLog
	for rows.Next() {
		var l domain.HabitLog
		if err := rows.Scan(&l.HabitID, &l.LogDate, &l.Completed); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/tasks/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const taskColumns = `id, user_id, title, description, category, subcategory, due_date, due_time,
	repeat_rule, priority, status, quadrant, sort_order, scheduled_at, subtasks, note_slug,
	completed_at, created_at, updated_at`

func (r *PostgresRepository) Create(ctx context.Context, t *domain.Task) error {
	subtasksJSON, err := json.Marshal(t.Subtasks)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx, `
		INSERT INTO tasks (user_id, title, description, category, subcategory, due_date, due_time,
			repeat_rule, priority, status, quadrant, sort_order, scheduled_at, subtasks, note_slug)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`, t.UserID, t.Title, t.Description, t.Category, t.Subcategory, t.DueDate, t.DueTime,
		t.RepeatRule, t.Priority, t.Status, t.Quadrant, t.SortOrder, t.ScheduledAt, subtasksJSON, t.NoteSlug,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *PostgresRepository) Update(ctx context.Context, t *domain.Task) error {
	subtasksJSON, err := json.Marshal(t.Subtasks)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE tasks SET
			title = $2, description = $3, category = $4, subcategory = $5, due_date = $6, due_time = $7,
			repeat_rule = $8, priority = $9, status = $10, quadrant = $11, sort_order = $12,
			scheduled_at = $13, subtasks = $14, note_slug = $15, completed_at = $16, updated_at = now()
		WHERE id = $1
	`, t.ID, t.Title, t.Description, t.Category, t.Subcategory, t.DueDate, t.DueTime,
		t.RepeatRule, t.Priority, t.Status, t.Quadrant, t.SortOrder, t.ScheduledAt, subtasksJSON, t.NoteSlug, t.CompletedAt)
	return err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+taskColumns+` FROM tasks WHERE id = $1`, id)
	t, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return t, err
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string, filter domain.ListFilter) ([]domain.Task, error) {
	query := `SELECT ` + taskColumns + ` FROM tasks WHERE user_id = $1`
	args := []interface{}{userID}

	if filter.Status != nil {
		args = append(args, *filter.Status)
		query += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if filter.Quadrant != nil {
		args = append(args, *filter.Quadrant)
		query += fmt.Sprintf(" AND quadrant = $%d", len(args))
	}
	if filter.From != nil {
		args = append(args, *filter.From)
		query += fmt.Sprintf(" AND due_date >= $%d", len(args))
	}
	if filter.To != nil {
		args = append(args, *filter.To)
		query += fmt.Sprintf(" AND due_date <= $%d", len(args))
	}
	query += " ORDER BY sort_order, created_at"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) ReorderSortOrders(ctx context.Context, userID string, orderedIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, id := range orderedIDs {
		if _, err := tx.Exec(ctx, `UPDATE tasks SET sort_order = $2 WHERE id = $1 AND user_id = $3`, id, i, userID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanTask(row rowScanner) (*domain.Task, error) {
	t := &domain.Task{}
	var subtasksJSON []byte
	var quadrant *string
	err := row.Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.Category, &t.Subcategory, &t.DueDate, &t.DueTime,
		&t.RepeatRule, &t.Priority, &t.Status, &quadrant, &t.SortOrder, &t.ScheduledAt, &subtasksJSON, &t.NoteSlug,
		&t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if quadrant != nil {
		q := domain.Quadrant(*quadrant)
		t.Quadrant = &q
	}
	if len(subtasksJSON) > 0 {
		if err := json.Unmarshal(subtasksJSON, &t.Subtasks); err != nil {
			return nil, err
		}
	}
	return t, nil
}

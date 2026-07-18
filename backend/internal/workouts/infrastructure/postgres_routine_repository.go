package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/workouts/domain"
)

type PostgresRoutineRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRoutineRepository(pool *pgxpool.Pool) *PostgresRoutineRepository {
	return &PostgresRoutineRepository{pool: pool}
}

func (r *PostgresRoutineRepository) CreateRoutine(ctx context.Context, routine *domain.Routine) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO workout_routines (user_id, name, exercises)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, routine.UserID, routine.Name, routine.Exercises).Scan(&routine.ID, &routine.CreatedAt)
}

func (r *PostgresRoutineRepository) ListRoutinesByUser(ctx context.Context, userID string) ([]domain.Routine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, exercises, created_at
		FROM workout_routines WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []domain.Routine
	for rows.Next() {
		var routine domain.Routine
		if err := rows.Scan(&routine.ID, &routine.UserID, &routine.Name, &routine.Exercises, &routine.CreatedAt); err != nil {
			return nil, err
		}
		routines = append(routines, routine)
	}
	return routines, rows.Err()
}

func (r *PostgresRoutineRepository) DeleteRoutine(ctx context.Context, userID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workout_routines WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

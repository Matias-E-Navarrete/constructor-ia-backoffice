package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/workouts/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, s *domain.Session) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO workout_sessions (user_id, session_date, notes, note_slug)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, s.UserID, s.SessionDate, s.Notes, s.NoteSlug).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return err
	}

	for _, set := range s.Sets {
		if _, err := tx.Exec(ctx, `
			INSERT INTO workout_sets (session_id, exercise_name, set_number, reps, weight_kg)
			VALUES ($1, $2, $3, $4, $5)
		`, s.ID, set.ExerciseName, set.SetNumber, set.Reps, set.WeightKg); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Session, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, session_date, notes, note_slug, created_at
		FROM workout_sessions WHERE user_id = $1 ORDER BY session_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.Session
	for rows.Next() {
		s := domain.Session{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.SessionDate, &s.Notes, &s.NoteSlug, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range sessions {
		sets, err := r.listSets(ctx, sessions[i].ID)
		if err != nil {
			return nil, err
		}
		sessions[i].Sets = sets
	}
	return sessions, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	s := &domain.Session{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, session_date, notes, note_slug, created_at
		FROM workout_sessions WHERE id = $1
	`, id).Scan(&s.ID, &s.UserID, &s.SessionDate, &s.Notes, &s.NoteSlug, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	sets, err := r.listSets(ctx, id)
	if err != nil {
		return nil, err
	}
	s.Sets = sets
	return s, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workout_sessions WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) ProgressForExercise(ctx context.Context, userID, exerciseName string) ([]domain.ExerciseEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ws.session_date, s.set_number, s.reps, s.weight_kg
		FROM workout_sets s
		JOIN workout_sessions ws ON ws.id = s.session_id
		WHERE ws.user_id = $1 AND s.exercise_name = $2
		ORDER BY ws.session_date
	`, userID, exerciseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.ExerciseEntry
	for rows.Next() {
		var e domain.ExerciseEntry
		if err := rows.Scan(&e.SessionDate, &e.SetNumber, &e.Reps, &e.WeightKg); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *PostgresRepository) AllPersonalRecords(ctx context.Context, userID string) ([]domain.ExercisePR, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.exercise_name, MAX(s.weight_kg)
		FROM workout_sets s
		JOIN workout_sessions ws ON ws.id = s.session_id
		WHERE ws.user_id = $1
		GROUP BY s.exercise_name
		ORDER BY s.exercise_name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []domain.ExercisePR
	for rows.Next() {
		var pr domain.ExercisePR
		if err := rows.Scan(&pr.ExerciseName, &pr.PersonalRecordKg); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}

func (r *PostgresRepository) listSets(ctx context.Context, sessionID string) ([]domain.Set, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT exercise_name, set_number, reps, weight_kg
		FROM workout_sets WHERE session_id = $1 ORDER BY set_number
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []domain.Set
	for rows.Next() {
		var set domain.Set
		if err := rows.Scan(&set.ExerciseName, &set.SetNumber, &set.Reps, &set.WeightKg); err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, rows.Err()
}

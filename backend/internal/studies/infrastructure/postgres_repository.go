package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/studies/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateSubject(ctx context.Context, s *domain.Subject) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO study_subjects (user_id, name, note_slug)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, s.UserID, s.Name, s.NoteSlug).Scan(&s.ID, &s.CreatedAt)
}

func (r *PostgresRepository) ListSubjectsByUser(ctx context.Context, userID string) ([]domain.Subject, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, note_slug, created_at
		FROM study_subjects WHERE user_id = $1 ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.NoteSlug, &s.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, rows.Err()
}

func (r *PostgresRepository) FindSubjectByID(ctx context.Context, id string) (*domain.Subject, error) {
	s := &domain.Subject{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, note_slug, created_at
		FROM study_subjects WHERE id = $1
	`, id).Scan(&s.ID, &s.UserID, &s.Name, &s.NoteSlug, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *PostgresRepository) DeleteSubject(ctx context.Context, userID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM study_subjects WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *PostgresRepository) LogSession(ctx context.Context, s *domain.StudySession) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO study_sessions (user_id, subject_id, session_date, duration_minutes, topic)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, s.UserID, s.SubjectID, s.SessionDate, s.DurationMinutes, s.Topic).Scan(&s.ID, &s.CreatedAt)
}

func (r *PostgresRepository) ListSessionsBySubject(ctx context.Context, subjectID string) ([]domain.StudySession, error) {
	return r.querySessions(ctx, `
		SELECT id, user_id, subject_id, session_date, duration_minutes, topic, created_at
		FROM study_sessions WHERE subject_id = $1 ORDER BY session_date DESC
	`, subjectID)
}

func (r *PostgresRepository) ListSessionsByUser(ctx context.Context, userID string) ([]domain.StudySession, error) {
	return r.querySessions(ctx, `
		SELECT id, user_id, subject_id, session_date, duration_minutes, topic, created_at
		FROM study_sessions WHERE user_id = $1 ORDER BY session_date DESC
	`, userID)
}

func (r *PostgresRepository) querySessions(ctx context.Context, query string, arg string) ([]domain.StudySession, error) {
	rows, err := r.pool.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.StudySession
	for rows.Next() {
		var s domain.StudySession
		if err := rows.Scan(&s.ID, &s.UserID, &s.SubjectID, &s.SessionDate, &s.DurationMinutes, &s.Topic, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

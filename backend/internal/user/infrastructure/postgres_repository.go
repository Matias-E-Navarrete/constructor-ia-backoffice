package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/user/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, u *domain.User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, plan, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, u.Email, u.PasswordHash, u.Plan, u.Role).Scan(&u.ID, &u.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrEmailTaken
	}
	return err
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.scanOne(r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, plan, role, created_at FROM users WHERE email = $1
	`, email))
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.scanOne(r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, plan, role, created_at FROM users WHERE id = $1
	`, id))
}

func (r *PostgresRepository) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, email, password_hash, plan, role, created_at FROM users ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Plan, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *PostgresRepository) UpdatePlan(ctx context.Context, id string, plan domain.Plan) (*domain.User, error) {
	return r.scanOne(r.pool.QueryRow(ctx, `
		UPDATE users SET plan = $2 WHERE id = $1
		RETURNING id, email, password_hash, plan, role, created_at
	`, id, plan))
}

func (r *PostgresRepository) UpdateRole(ctx context.Context, id string, role domain.Role) (*domain.User, error) {
	return r.scanOne(r.pool.QueryRow(ctx, `
		UPDATE users SET role = $2 WHERE id = $1
		RETURNING id, email, password_hash, plan, role, created_at
	`, id, role))
}

func (r *PostgresRepository) scanOne(row pgx.Row) (*domain.User, error) {
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Plan, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

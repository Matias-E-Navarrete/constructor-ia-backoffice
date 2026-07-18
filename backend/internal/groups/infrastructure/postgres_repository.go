package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"rimu/backend/internal/groups/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, g *domain.Group, ownerRole domain.MemberRole) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO groups (name, kind, owner_user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, g.Name, g.Kind, g.OwnerUserID).Scan(&g.ID, &g.CreatedAt)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)
	`, g.ID, g.OwnerUserID, ownerRole); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) AddMember(ctx context.Context, groupID, userID string, role domain.MemberRole) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3)
	`, groupID, userID, role)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrAlreadyMember
	}
	return err
}

func (r *PostgresRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM group_members WHERE group_id = $1 AND user_id = $2
	`, groupID, userID)
	return err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Group, []domain.Member, error) {
	g := &domain.Group{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, kind, owner_user_id, created_at FROM groups WHERE id = $1
	`, id).Scan(&g.ID, &g.Name, &g.Kind, &g.OwnerUserID, &g.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT group_id, user_id, role, joined_at FROM group_members WHERE group_id = $1
	`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, nil, err
		}
		members = append(members, m)
	}
	return g, members, rows.Err()
}

func (r *PostgresRepository) ListForUser(ctx context.Context, userID string) ([]domain.Membership, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.name, g.kind, g.owner_user_id, g.created_at, gm.role
		FROM groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE gm.user_id = $1
		ORDER BY g.created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []domain.Membership
	for rows.Next() {
		var m domain.Membership
		if err := rows.Scan(&m.Group.ID, &m.Group.Name, &m.Group.Kind, &m.Group.OwnerUserID, &m.Group.CreatedAt, &m.Role); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, rows.Err()
}

func (r *PostgresRepository) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)
	`, groupID, userID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) HasCoachAccess(ctx context.Context, coachUserID, clientUserID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM groups g
			JOIN group_members coach ON coach.group_id = g.id AND coach.user_id = $1 AND coach.role = 'coach'
			JOIN group_members client ON client.group_id = g.id AND client.user_id = $2 AND client.role = 'client'
			WHERE g.kind = 'coaching'
		)
	`, coachUserID, clientUserID).Scan(&exists)
	return exists, err
}

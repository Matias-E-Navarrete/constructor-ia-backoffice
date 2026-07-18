package application

import (
	"context"
	"errors"
	"strings"

	"rimu/backend/internal/platform/auth"
	"rimu/backend/internal/user/domain"
)

var ErrWeakPassword = errors.New("user: password must be at least 8 characters")

type RegisterUser struct {
	Repo        domain.Repository
	IsAdminEmail func(email string) bool
}

type RegisterUserInput struct {
	Email    string
	Password string
}

func (uc *RegisterUser) Execute(ctx context.Context, in RegisterUserInput) (*domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, errors.New("user: invalid email")
	}
	if len(in.Password) < 8 {
		return nil, ErrWeakPassword
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	role := domain.RoleUser
	if uc.IsAdminEmail != nil && uc.IsAdminEmail(email) {
		role = domain.RoleAdmin
	}

	u := &domain.User{
		Email:        email,
		PasswordHash: hash,
		Plan:         domain.PlanFree,
		Role:         role,
	}
	if err := uc.Repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

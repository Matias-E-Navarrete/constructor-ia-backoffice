package application

import (
	"context"
	"errors"

	"rimu/backend/internal/platform/auth"
	"rimu/backend/internal/user/domain"
)

var ErrInvalidCredentials = errors.New("user: invalid email or password")

type AuthenticateUser struct {
	Repo      domain.Repository
	JWTSecret string
}

type AuthenticateUserInput struct {
	Email    string
	Password string
}

type AuthenticateUserOutput struct {
	Token string
	User  *domain.User
}

func (uc *AuthenticateUser) Execute(ctx context.Context, in AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	u, err := uc.Repo.FindByEmail(ctx, in.Email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.CheckPassword(u.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.IssueToken(uc.JWTSecret, u.ID)
	if err != nil {
		return nil, err
	}
	return &AuthenticateUserOutput{Token: token, User: u}, nil
}

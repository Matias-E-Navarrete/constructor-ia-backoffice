package application

import (
	"context"

	"rimu/backend/internal/finance/domain"
)

type CreateAccountInput struct {
	UserID         string
	Name           string
	Type           string
	Currency       string
	InitialBalance float64
}

type CreateAccount struct {
	Repo domain.AccountRepository
}

func (uc *CreateAccount) Execute(ctx context.Context, in CreateAccountInput) (*domain.Account, error) {
	a := &domain.Account{
		UserID: in.UserID, Name: in.Name, Type: in.Type, Currency: in.Currency,
		InitialBalance: in.InitialBalance, Active: true,
	}
	if err := uc.Repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

type ListAccounts struct {
	Repo domain.AccountRepository
}

func (uc *ListAccounts) Execute(ctx context.Context, userID string) ([]domain.Account, error) {
	return uc.Repo.ListByUser(ctx, userID)
}

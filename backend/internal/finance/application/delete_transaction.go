package application

import (
	"context"
	"errors"

	"rimu/backend/internal/finance/domain"
)

var ErrNotOwner = errors.New("finance: you do not own this transaction")

type DeleteTransaction struct {
	Repo domain.Repository
}

func (uc *DeleteTransaction) Execute(ctx context.Context, userID, txID string) error {
	t, err := uc.Repo.FindByID(ctx, txID)
	if err != nil {
		return err
	}
	if t.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.Delete(ctx, txID)
}

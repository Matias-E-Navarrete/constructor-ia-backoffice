package application

import (
	"context"
	"errors"
	"time"

	groupsdomain "rimu/backend/internal/groups/domain"
	"rimu/backend/internal/finance/domain"
)

var ErrNotGroupMember = errors.New("finance: you are not a member of that group")

type RecordTransactionInput struct {
	UserID      string
	Type        domain.TxType
	AmountValue float64
	Category    string
	Description string
	TxDate      time.Time
	NoteSlug    *string
	GroupID     *string
}

type RecordTransaction struct {
	Repo   domain.Repository
	Groups groupsdomain.Repository
}

func (uc *RecordTransaction) Execute(ctx context.Context, in RecordTransactionInput) (*domain.Transaction, error) {
	amount, err := domain.NewAmount(in.AmountValue)
	if err != nil {
		return nil, err
	}

	if in.GroupID != nil {
		isMember, err := uc.Groups.IsMember(ctx, *in.GroupID, in.UserID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, ErrNotGroupMember
		}
	}

	t := &domain.Transaction{
		UserID:      in.UserID,
		Type:        in.Type,
		Amount:      amount,
		Category:    in.Category,
		Description: in.Description,
		TxDate:      in.TxDate,
		NoteSlug:    in.NoteSlug,
		GroupID:     in.GroupID,
	}
	if err := uc.Repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

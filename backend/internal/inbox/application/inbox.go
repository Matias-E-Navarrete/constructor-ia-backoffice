package application

import (
	"context"
	"errors"

	"rimu/backend/internal/inbox/domain"
)

var ErrNotOwner = errors.New("inbox: you do not own this item")

type CaptureItem struct {
	Repo domain.Repository
}

func (uc *CaptureItem) Execute(ctx context.Context, userID, content string) (*domain.Item, error) {
	item := &domain.Item{UserID: userID, Content: content}
	if err := uc.Repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

type ListInbox struct {
	Repo domain.Repository
}

func (uc *ListInbox) Execute(ctx context.Context, userID string) ([]domain.Item, error) {
	return uc.Repo.ListByUser(ctx, userID)
}

type SetPinned struct {
	Repo domain.Repository
}

func (uc *SetPinned) Execute(ctx context.Context, userID, itemID string, pinned bool) error {
	item, err := uc.Repo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.SetPinned(ctx, itemID, pinned)
}

type DeleteItem struct {
	Repo domain.Repository
}

func (uc *DeleteItem) Execute(ctx context.Context, userID, itemID string) error {
	item, err := uc.Repo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrNotOwner
	}
	return uc.Repo.Delete(ctx, itemID)
}

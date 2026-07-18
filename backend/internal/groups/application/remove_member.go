package application

import (
	"context"

	"rimu/backend/internal/groups/domain"
)

type RemoveMember struct {
	Repo domain.Repository
}

type RemoveMemberInput struct {
	RequesterUserID string
	GroupID         string
	TargetUserID    string
}

func (uc *RemoveMember) Execute(ctx context.Context, in RemoveMemberInput) error {
	_, members, err := uc.Repo.FindByID(ctx, in.GroupID)
	if err != nil {
		return err
	}
	// A member can always remove themselves; otherwise only the owner/coach can.
	if in.RequesterUserID != in.TargetUserID && !isOwnerOrCoach(members, in.RequesterUserID) {
		return ErrNotAuthorized
	}
	return uc.Repo.RemoveMember(ctx, in.GroupID, in.TargetUserID)
}

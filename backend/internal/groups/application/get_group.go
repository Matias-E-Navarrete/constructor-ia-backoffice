package application

import (
	"context"
	"errors"

	"rimu/backend/internal/groups/domain"
)

var ErrNotMember = errors.New("groups: you are not a member of this group")

type GetGroup struct {
	Repo domain.Repository
}

type GetGroupOutput struct {
	Group   domain.Group
	Members []domain.Member
}

func (uc *GetGroup) Execute(ctx context.Context, requesterUserID, groupID string) (*GetGroupOutput, error) {
	group, members, err := uc.Repo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if !isMemberOf(members, requesterUserID) {
		return nil, ErrNotMember
	}
	return &GetGroupOutput{Group: *group, Members: members}, nil
}

func isMemberOf(members []domain.Member, userID string) bool {
	for _, m := range members {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

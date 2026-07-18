package application

import (
	"context"
	"errors"

	"rimu/backend/internal/groups/domain"
	userdomain "rimu/backend/internal/user/domain"
)

var ErrNotAuthorized = errors.New("groups: only the owner/coach can invite members")

type InviteMember struct {
	Groups domain.Repository
	Users  userdomain.Repository
}

type InviteMemberInput struct {
	RequesterUserID string
	GroupID         string
	InviteeEmail    string
}

func (uc *InviteMember) Execute(ctx context.Context, in InviteMemberInput) (*domain.Member, error) {
	group, members, err := uc.Groups.FindByID(ctx, in.GroupID)
	if err != nil {
		return nil, err
	}
	if !isOwnerOrCoach(members, in.RequesterUserID) {
		return nil, ErrNotAuthorized
	}

	invitee, err := uc.Users.FindByEmail(ctx, in.InviteeEmail)
	if err != nil {
		return nil, err
	}

	role := domain.RoleMember
	if group.Kind == domain.KindCoaching {
		role = domain.RoleClient
	}

	if err := uc.Groups.AddMember(ctx, group.ID, invitee.ID, role); err != nil {
		return nil, err
	}
	return &domain.Member{GroupID: group.ID, UserID: invitee.ID, Role: role}, nil
}

func isOwnerOrCoach(members []domain.Member, userID string) bool {
	for _, m := range members {
		if m.UserID == userID && (m.Role == domain.RoleOwner || m.Role == domain.RoleCoach) {
			return true
		}
	}
	return false
}

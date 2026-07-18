package domain

import (
	"context"
	"errors"
	"time"
)

type Kind string

const (
	KindFamily   Kind = "family"
	KindCoaching Kind = "coaching"
)

type MemberRole string

const (
	RoleOwner  MemberRole = "owner"
	RoleMember MemberRole = "member"
	RoleCoach  MemberRole = "coach"
	RoleClient MemberRole = "client"
)

var ErrNotFound = errors.New("groups: not found")
var ErrAlreadyMember = errors.New("groups: user is already a member")

type Group struct {
	ID          string
	Name        string
	Kind        Kind
	OwnerUserID string
	CreatedAt   time.Time
}

type Member struct {
	GroupID  string
	UserID   string
	Role     MemberRole
	JoinedAt time.Time
}

// Membership is a read-model row: "the groups I belong to, and my role in each."
type Membership struct {
	Group Group
	Role  MemberRole
}

type Repository interface {
	Create(ctx context.Context, g *Group, ownerRole MemberRole) error
	AddMember(ctx context.Context, groupID, userID string, role MemberRole) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	FindByID(ctx context.Context, id string) (*Group, []Member, error)
	ListForUser(ctx context.Context, userID string) ([]Membership, error)
	IsMember(ctx context.Context, groupID, userID string) (bool, error)
	// HasCoachAccess reports whether coachUserID has an active 'coach' role
	// over clientUserID within some 'coaching' group. Used by habits/workouts
	// to authorize a coach reading a client's data.
	HasCoachAccess(ctx context.Context, coachUserID, clientUserID string) (bool, error)
}

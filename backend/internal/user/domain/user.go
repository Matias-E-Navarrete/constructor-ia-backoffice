package domain

import (
	"context"
	"errors"
	"time"
)

type Plan string

const (
	PlanFree Plan = "free"
	PlanPro  Plan = "pro"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var ErrEmailTaken = errors.New("user: email already registered")
var ErrNotFound = errors.New("user: not found")

type User struct {
	ID               string
	Email            string
	PasswordHash     string
	Plan             Plan
	Role             Role
	PlannerStartHour int
	PlannerEndHour   int
	CreatedAt        time.Time
}

func (u *User) IsPro() bool {
	return u.Plan == PlanPro
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Upgrade flips the account to the pro plan. In this MVP it's called
// directly by a mock endpoint; a real payment provider would call it from a
// webhook instead.
func (u *User) Upgrade() {
	u.Plan = PlanPro
}

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	UpdatePlan(ctx context.Context, id string, plan Plan) (*User, error)
	UpdateRole(ctx context.Context, id string, role Role) (*User, error)
	UpdatePlannerHours(ctx context.Context, id string, startHour, endHour int) (*User, error)
}

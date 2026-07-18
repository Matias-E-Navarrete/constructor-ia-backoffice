package http

import "rimu/backend/internal/user/domain"

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
	Role  string `json:"role"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Plan: string(u.Plan), Role: string(u.Role)}
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

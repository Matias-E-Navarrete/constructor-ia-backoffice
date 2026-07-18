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

type updatePlannerHoursRequest struct {
	StartHour int `json:"start_hour"`
	EndHour   int `json:"end_hour"`
}

type userResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Plan             string `json:"plan"`
	Role             string `json:"role"`
	PlannerStartHour int    `json:"planner_start_hour"`
	PlannerEndHour   int    `json:"planner_end_hour"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID: u.ID, Email: u.Email, Plan: string(u.Plan), Role: string(u.Role),
		PlannerStartHour: u.PlannerStartHour, PlannerEndHour: u.PlannerEndHour,
	}
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

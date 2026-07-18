package http

import (
	admindomain "rimu/backend/internal/admin/domain"
	ffdomain "rimu/backend/internal/platform/featureflags/domain"
	userdomain "rimu/backend/internal/user/domain"
)

type createRoadmapItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Kind        string `json:"kind,omitempty"`
}

type updateRoadmapStatusRequest struct {
	Status string `json:"status"`
}

type roadmapItemResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Status      string `json:"status"`
}

func toRoadmapItemResponse(it admindomain.RoadmapItem) roadmapItemResponse {
	return roadmapItemResponse{ID: it.ID, Title: it.Title, Description: it.Description, Kind: string(it.Kind), Status: string(it.Status)}
}

type toggleFlagRequest struct {
	Enabled bool `json:"enabled"`
}

type flagResponse struct {
	Key         string `json:"key"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

func toFlagResponse(f ffdomain.Flag) flagResponse {
	return flagResponse{Key: f.Key, Enabled: f.Enabled, Description: f.Description}
}

type setUserPlanRequest struct {
	Plan string `json:"plan"`
}

type setUserRoleRequest struct {
	Role string `json:"role"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
	Role  string `json:"role"`
}

func toUserResponse(u *userdomain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Plan: string(u.Plan), Role: string(u.Role)}
}

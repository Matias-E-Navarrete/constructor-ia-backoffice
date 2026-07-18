package http

import "rimu/backend/internal/groups/domain"

type createGroupRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type inviteMemberRequest struct {
	Email string `json:"email"`
}

type groupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	OwnerUserID string `json:"owner_user_id"`
}

func toGroupResponse(g domain.Group) groupResponse {
	return groupResponse{ID: g.ID, Name: g.Name, Kind: string(g.Kind), OwnerUserID: g.OwnerUserID}
}

type membershipResponse struct {
	Group groupResponse `json:"group"`
	Role  string        `json:"role"`
}

type memberResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type groupDetailResponse struct {
	Group   groupResponse    `json:"group"`
	Members []memberResponse `json:"members"`
}

package http

import "rimu/backend/internal/inbox/domain"

type captureRequest struct {
	Content string `json:"content"`
}

type setPinnedRequest struct {
	Pinned bool `json:"pinned"`
}

type itemResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Pinned  bool   `json:"pinned"`
}

func toItemResponse(it domain.Item) itemResponse {
	return itemResponse{ID: it.ID, Content: it.Content, Pinned: it.Pinned}
}

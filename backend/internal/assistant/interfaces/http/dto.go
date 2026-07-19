package http

import "rimu/backend/internal/assistant/domain"

type sendMessageRequest struct {
	Content string `json:"content"`
}

type messageResponse struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func toMessageResponse(m domain.Message) messageResponse {
	return messageResponse{
		ID:        m.ID,
		Role:      string(m.Role),
		Content:   m.Content,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

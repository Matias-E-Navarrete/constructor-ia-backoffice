package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"rimu/backend/internal/assistant/domain"
)

const systemPrompt = `Sos el asistente de rimu, un sistema personal que combina Tareas ` +
	`(Eisenhower/Kanban), Planner y Modo Foco, Hábitos, Entrenamientos, Estudios, Wallet ` +
	`y Notas estilo Obsidian. Ayudás al usuario a organizarse, priorizar, mantener rachas ` +
	`de hábitos y estudio, y a pensar su semana. Respondé siempre en español, en tono ` +
	`cercano y directo, con respuestas cortas salvo que pidan detalle.`

// AnthropicClient is a minimal wrapper around the Messages API — no SDK
// dependency, matching the project's "plain net/http" convention (same
// spirit as the PDF export using fpdf directly).
type AnthropicClient struct {
	APIKey  string
	Model   string
	BaseURL string
	Client  *http.Client
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		APIKey:  apiKey,
		Model:   "claude-sonnet-5",
		BaseURL: "https://api.anthropic.com/v1/messages",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *AnthropicClient) Complete(ctx context.Context, history []domain.Message) (string, error) {
	if c.APIKey == "" {
		return "", domain.ErrNotConfigured
	}

	messages := make([]anthropicMessage, 0, len(history))
	for _, m := range history {
		messages = append(messages, anthropicMessage{Role: string(m.Role), Content: m.Content})
	}

	body, err := json.Marshal(anthropicRequest{
		Model:     c.Model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages:  messages,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("assistant: decoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return "", fmt.Errorf("assistant: provider error: %s", parsed.Error.Message)
		}
		return "", errors.New("assistant: provider returned a non-200 status")
	}
	if len(parsed.Content) == 0 {
		return "", errors.New("assistant: provider returned no content")
	}
	return parsed.Content[0].Text, nil
}

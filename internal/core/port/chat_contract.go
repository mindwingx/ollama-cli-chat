package port

import "ollama-cli/internal/core/domain"

type IChatService interface {
	Handshake() (bool, error)
	SendChatMessage(model, role, msg string, stream bool) (domain.Message, error)
	GetModelsList() (domain.Models, error)
	PullModel(model string) error
	DeleteModel(model string) error
}

package chat

import (
	"chat/internal/app/models"
	"context"
)

// Usecase - port (первичный)
type Usecase interface {
	CreateChat(ctx context.Context, dto *CreateChatDTO) (*models.Chat, error)
}

type (
	// ChatRepository - port (вторичный)
	ChatRepository interface {
		CreateChat(ctx context.Context, order *models.Chat) error
	}
)

// Deps -
type Deps struct {
	// Adapters
	ChatRepository ChatRepository
}

type usecase struct {
	Deps
}

func NewUsecase(d Deps) Usecase {
	return &usecase{Deps: d}
}

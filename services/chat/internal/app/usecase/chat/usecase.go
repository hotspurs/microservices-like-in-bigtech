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
	// User - port (вторичный)
	User interface {
		Check(ctx context.Context, users []models.UserID) error
	}
	// ChatRepository - port (вторичный)
	ChatRepository interface {
		CreateChat(ctx context.Context, order *models.Chat) error
	}
)

// Deps -
type Deps struct {
	// Adapters
	ChatRepository ChatRepository
	User           User
}

type usecase struct {
	Deps
}

func NewUsecase(d Deps) Usecase {
	return &usecase{Deps: d}
}

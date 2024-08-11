package chat

import (
	"chat/internal/app/models"
	"context"
)

type Usecase interface {
	CreateChat(ctx context.Context, dto *CreateChatDTO) (*models.Chat, error)
}

// Deps -
type Deps struct {
	// Adapters
}

type usecase struct {
	Deps
}

func NewUsecase(d Deps) Usecase {
	return &usecase{Deps: d}
}

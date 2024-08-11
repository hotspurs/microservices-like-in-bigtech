package chat

import (
	"chat/internal/app/models"
	"context"
)

func (uc *usecase) CreateChat(ctx context.Context, chatInfo *CreateChatDTO) (*models.Chat, error) {
	return &models.Chat{ID: 1}, nil
}

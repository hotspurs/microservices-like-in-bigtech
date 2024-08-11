package chat

import (
	"chat/internal/app/models"
	"context"
)

func (uc *usecase) CreateChat(ctx context.Context, chatInfo *CreateChatDTO) (*models.Chat, error) {
	if err := uc.Deps.User.Check(ctx, chatInfo.UserIds); err != nil {
		return nil, err // здесь мы возвращаем ошибку err
	}

	chat := models.NewChat()

	if err := uc.Deps.ChatRepository.CreateChat(ctx, chat); err != nil {
		return nil, err
	}

	return &models.Chat{ID: 1}, nil
}

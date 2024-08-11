package chat

import "chat/internal/app/models"

type (
	CreateChatDTO struct {
		UserIds []models.UserID
	}
)

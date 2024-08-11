package models

import "github.com/google/uuid"

type ChatID uint64

type Chat struct {
	ID ChatID
}

func NewChat() *Chat {
	return &Chat{
		ID: ChatID(uuid.New().ID()),
	}
}

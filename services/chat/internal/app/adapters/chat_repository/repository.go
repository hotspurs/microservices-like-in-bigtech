package chatrepository

import "chat/internal/app/usecase/chat"

type Repository struct {
	// postgres/mongodb/redis/...
}

var (
	_ chat.ChatRepository = (*Repository)(nil)
)

func NewRepository( /**/ ) *Repository {
	return &Repository{}
}

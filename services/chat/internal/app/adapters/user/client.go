package user

import "chat/internal/app/usecase/chat"

type Client struct {
	// http/grpc/...
}

var (
	_ chat.User = (*Client)(nil)
)

func NewClient( /**/ ) *Client {
	return &Client{}
}

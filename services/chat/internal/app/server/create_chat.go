package server

import (
	"chat/internal/app/usecase/chat"
	pb "chat/pkg/api/chat"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	// 1. validation
	if err := validateCreateChatRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// 2. convert delivery models to domain models/DTO
	chatDto := newChatFromPbCreateChatRequest(req)

	// 3. call usecase
	newChat, err := s.ChatUsecase.CreateChat(ctx, chatDto)
	if err != nil {
		return nil, err // обработается на уровне middleware
	}

	// 4. convert domain models/DTO to delivery models
	response := &pb.CreateChatResponse{
		Id: uint64(newChat.ID),
	}

	// 5. return result
	return &pb.CreateChatResponse{
		Id: response.Id,
	}, nil
}

func validateCreateChatRequest(_ *pb.CreateChatRequest) error {
	//
	return nil
}

func newChatFromPbCreateChatRequest(req *pb.CreateChatRequest) *chat.CreateChatDTO {
	return &chat.CreateChatDTO{
		UserIds: req.GetUserIds(),
	}
}

package main

import (
	chatrepository "chat/internal/app/adapters/chat_repository"
	"chat/internal/app/adapters/user"
	"chat/internal/app/server"
	"chat/internal/app/usecase/chat"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// =========================
	// adapters
	// =========================

	// repository
	chatRepository := chatrepository.NewRepository()

	// services
	userAdapter := user.NewClient()

	// =========================
	// usecases
	// =========================

	chatUsecase := chat.NewUsecase(chat.Deps{
		chatRepository,
		userAdapter,
	})

	// =========================
	// delivery
	// =========================

	config := server.Config{
		GRPCPort:               ":8082",
		GRPCGatewayPort:        ":8080",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{},
	}

	srv, err := server.New(ctx, config, server.Deps{
		ChatUsecase: chatUsecase,
		// Dependency injection
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}

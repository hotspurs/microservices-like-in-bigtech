package main

import (
	pb "chat/pkg/api/chat"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"net/http"
	"sync"
)

type server struct {
	pb.UnimplementedChatServiceServer
	mx sync.RWMutex
}

func NewServer() (*server, error) {
	return &server{}, nil
}

func (s *server) SendMessage(_ context.Context, _ *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return &pb.SendMessageResponse{
		Message: &pb.ChatMessage{
			Id:        1,
			Text:      "text",
			UserId:    1,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server, err := NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		grpcServer := grpc.NewServer()
		pb.RegisterChatServiceServer(grpcServer, server)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Register gRPC server endpoint
		// Note: Make sure the gRPC server is running properly and accessible
		mux := runtime.NewServeMux()
		if err = pb.RegisterChatServiceHandlerServer(ctx, mux, server); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		e := echo.New()
		e.GET("/live", func(c echo.Context) error {
			return c.String(http.StatusOK, "")
		})

		e.GET("/ready", func(c echo.Context) error {
			return c.String(http.StatusOK, "")
		})
		e.GET("/auth", func(c echo.Context) error {
			res, err := http.Get("http://auth:8080")
			if err != nil {
			}
			fmt.Println(res)
			return c.String(res.StatusCode, "")
		})
		e.Static("/docs", "./swagger")
		e.Any("/*", echo.WrapHandler(mux))

		e.Logger.Fatal(e.Start(":8080"))
	}()
	wg.Wait()
}

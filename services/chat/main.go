package main

import (
	pb "chat/pkg/api/chat"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

		// SaveNote:
		// grpc_cli call --json_input --json_output localhost:8082 NotesService/SaveNote '{"info":{"title":"my note","content":"my note content"}}'
		// ListNotes:
		// grpc_cli call --json_input --json_output localhost:8082 NotesService/ListNotes '{}'
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
		httpServer := &http.Server{Handler: mux}

		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Printf("server listening at %v", lis.Addr())
		if err := httpServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		e := echo.New()
		e.Static("/docs", "./swagger")
		e.Logger.Fatal(e.Start(":8083"))
	}()

	wg.Wait()
	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, Chat!")
	//})
	//e.GET("/live", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "")
	//})
	//
	//e.GET("/ready", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "")
	//})
	//e.GET("/auth", func(c echo.Context) error {
	//	res, err := http.Get("http://auth:8080")
	//	if err != nil {
	//	}
	//	fmt.Println(res)
	//	return c.String(res.StatusCode, "")
	//})
	//e.Logger.Fatal(e.Start(":8080"))
}

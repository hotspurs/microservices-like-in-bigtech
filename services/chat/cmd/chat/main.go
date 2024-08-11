package main

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pb "chat/pkg/api/chat"
	"context"
	"errors"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	mx        sync.RWMutex
	validator *protovalidate.Validator
}

func NewServer() (*server, error) {
	srv := &server{}
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.CreateChatRequest{},
			&pb.SendMessageRequest{},
			&pb.GetMessagesRequest{},
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	srv.validator = validator

	return srv, nil
}

func protovalidateVialationsToGoogleViolations(vs []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldPath,
			Description: v.Message,
		}
	}
	return res
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.Violations),
	}
}

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

func (s *server) SendMessage(_ context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.SendMessageResponse{
		Message: &pb.ChatMessage{
			Id:        1,
			Text:      "text",
			UserId:    1,
			Timestamp: timestamppb.Now(),
		},
	}, nil
}

func (s *server) GetMessages(_ context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.GetMessagesResponse{
		Items: make([]*pb.ChatMessage, 0),
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
		e.Static("/docs", "swagger")
		e.Any("/*", echo.WrapHandler(mux))

		e.Logger.Fatal(e.Start(":8080"))
	}()
	wg.Wait()
}

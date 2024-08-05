package main

import (
	pb "chat/pkg/api/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func NewServer() *server {
	return &server{}
}

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, NewServer())

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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

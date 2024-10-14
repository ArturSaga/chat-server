package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserApiServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserApiServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateChat - публичный метод, который создает новый чат.
func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	log.Printf("Context: %+v", ctx)
	for _, value := range req.Usernames {
		log.Printf("User id: %+v", value)
	}

	return &desc.CreateChatResponse{}, nil
}

// SendMessage - публичный метод, который отправляет сообщение в чат.
func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Context: %+v", ctx)
	log.Printf("User id: %+v", req)

	return &emptypb.Empty{}, nil
}

// DeleteChat - публичный метод, который удаляет чат.
func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	log.Printf("Context: %+v", ctx)
	log.Printf("User id: %+d", req.GetId())
	return &emptypb.Empty{}, nil
}

package chat

import (
	"context"
	"fmt"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	converter "github.com/ArturSaga/chat-server/internal/convertor"
)

// CreateChat - публичный метод, который создает пользователя.
func (i *Implementation) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	id, err := i.chatService.CreateChat(ctx, converter.ToChatFromDesc(req))

	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return nil, err
	}

	return &desc.CreateChatResponse{
		Id: id,
	}, nil
}

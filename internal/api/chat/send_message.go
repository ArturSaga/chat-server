package chat

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	converter "github.com/ArturSaga/chat-server/internal/convertor"
)

// SendMessage - публичный метод, который создает пользователя.
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	_, err := i.messageService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

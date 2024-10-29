package chat

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
)

// DeleteChat - публичный метод, который создает пользователя.
func (i *Implementation) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	_, err := i.chatService.DeleteChat(ctx, int64(req.Id))

	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

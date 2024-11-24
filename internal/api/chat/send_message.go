package chat

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	converter "github.com/ArturSaga/chat-server/internal/convertor"
	se "github.com/ArturSaga/chat-server/internal/service_error"
)

// SendMessage - публичный метод, который создает пользователя.
func (cs *ChatServer) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	_, err := cs.validateSendMessageRequest(req)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	_, err = cs.messageService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		fmt.Printf("failed to send message: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (cs *ChatServer) validateSendMessageRequest(req *desc.SendMessageRequest) (bool, error) {
	if req.From == "" || req.UserId == 0 || req.ChatId == 0 || req.Text == "" {
		return false, se.ErrRequireParam
	}

	return true, nil
}

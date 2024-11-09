package chat

import (
	"context"
	"fmt"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	converter "github.com/ArturSaga/chat-server/internal/convertor"
	se "github.com/ArturSaga/chat-server/internal/service_error"
)

// CreateChat - публичный метод, который создает пользователя.
func (cs *ChatServer) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	_, err := cs.validateCreateChatRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := cs.chatService.CreateChat(ctx, converter.ToChatFromDesc(req))
	if err != nil {
		fmt.Printf("failed to create chat: %v", err)
		return nil, err
	}

	return &desc.CreateChatResponse{
		Id: id,
	}, nil
}

func (cs *ChatServer) validateCreateChatRequest(req *desc.CreateChatRequest) (bool, error) {
	if req.ChatName == "" || req.UserIds == nil || req.Usernames == nil {
		return false, se.ErrRequireParam
	}

	if len(req.UserIds) != len(req.Usernames) {
		fmt.Println("count userIDs not match to count usernames")
		return false, se.ErrUserIDsNotMatchUserNames
	}

	return true, nil
}

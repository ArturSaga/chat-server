package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/model"
)

// ChatService - интерфейс, определяющий метод сервисного слоя чата
type ChatService interface {
	CreateChat(ctx context.Context, chat *model.Chat) (int64, error)
	DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error)
}

// MessageService - интерфейс, определяющий метод сервисного слоя сообщения
type MessageService interface {
	SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error)
}

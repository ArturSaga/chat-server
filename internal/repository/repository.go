package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/model"
)

// ChatRepository - интерфейс, определяющий методы репо слоя чата
type ChatRepository interface {
	CreateChat(ctx context.Context, chat *model.Chat) (int64, error)
	DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error)
}

// MessageRepository - интерфейс, определяющий методы репо слоя сообщений
type MessageRepository interface {
	SendMessage(ctx context.Context, messageInfo *model.Message) (*emptypb.Empty, error)
}

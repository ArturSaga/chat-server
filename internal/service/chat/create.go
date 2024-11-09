package chat

import (
	"context"

	"github.com/ArturSaga/chat-server/internal/model"
)

// CreateChat - публичный метод, для создания чата в слое сервиса.
func (s *chatService) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	id, err := s.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return 0, err
	}

	return id, nil
}

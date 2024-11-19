package chat

import (
	"github.com/ArturSaga/chat-server/internal/client/db"
	"github.com/ArturSaga/chat-server/internal/repository"
	"github.com/ArturSaga/chat-server/internal/service"
)

type chatService struct {
	chatRepo  repository.ChatRepository
	txManager db.TxManager
}

// NewChatService - публчиный метод, создающий сущность, для работы с сервисным слоем
func NewChatService(chatRepo repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &chatService{
		chatRepo:  chatRepo,
		txManager: txManager,
	}
}

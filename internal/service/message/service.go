package message

import (
	"github.com/ArturSaga/platform_common/pkg/db"

	"github.com/ArturSaga/chat-server/internal/repository"
	"github.com/ArturSaga/chat-server/internal/service"
)

type messageService struct {
	messageRepo repository.MessageRepository
	txManager   db.TxManager
}

// NewMessageService - публчиный метод, создающий сущность, для работы с сервисным слоем
func NewMessageService(messageRepo repository.MessageRepository, txManager db.TxManager) service.MessageService {
	return &messageService{
		messageRepo: messageRepo,
		txManager:   txManager,
	}
}

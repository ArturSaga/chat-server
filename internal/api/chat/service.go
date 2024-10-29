package chat

import (
	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/service"
)

// Implementation - сущность, которая ипмлементирует контракты
type Implementation struct {
	desc.UnimplementedChatApiServer
	chatService    service.ChatService
	messageService service.MessageService
}

// NewImplementation - публичный метод, реализует контракты
func NewImplementation(chatService service.ChatService, messageService service.MessageService) *Implementation {
	return &Implementation{
		chatService:    chatService,
		messageService: messageService,
	}
}

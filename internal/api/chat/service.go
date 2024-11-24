package chat

import (
	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/service"
)

// ChatServer - сущность, которая ипмлементирует контракты
type ChatServer struct {
	desc.UnimplementedChatApiServer
	chatService    service.ChatService
	messageService service.MessageService
}

// NewChatServer - публичный метод, реализует контракты
func NewChatServer(chatService service.ChatService, messageService service.MessageService) *ChatServer {
	return &ChatServer{
		chatService:    chatService,
		messageService: messageService,
	}
}

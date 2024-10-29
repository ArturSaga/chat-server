package convertor

import (
	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/model"
)

func ToMessageFromDesc(message *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		From:      message.From,
		ChatID:    message.ChatId,
		UserID:    message.UserId,
		Text:      message.Text,
		Timestamp: message.Timestamp.AsTime(),
	}
}
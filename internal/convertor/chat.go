package convertor

import (
	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/model"
)

func ToChatFromDesc(chat *desc.CreateChatRequest) *model.Chat {
	return &model.Chat{
		UserIDs:   chat.UserIds,
		UserNames: chat.Usernames,
		ChatName:  chat.ChatName,
	}
}

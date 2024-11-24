package message

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/model"
)

// SendMessage - публичный метод, для отправки сообщения в чат, в слое сервиса
func (s *messageService) SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error) {
	fmt.Println(message.Text)
	_, err := s.messageRepo.SendMessage(ctx, message)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

package message

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/model"
)

func (s *messageService) SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error) {
	fmt.Println(message.Text)
	_, err := s.messageRepo.SendMessage(ctx, message)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteChat - публичный метод, для удаления чата в слое сервиса
func (s *chatService) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	_, err := s.chatRepo.DeleteChat(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

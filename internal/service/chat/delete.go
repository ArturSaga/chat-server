package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *chatService) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	_, err := s.chatRepo.DeleteChat(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

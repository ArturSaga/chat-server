package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/api/chat"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/service"
	serviceMocks "github.com/ArturSaga/chat-server/internal/service/mocks"
	serviceErrors "github.com/ArturSaga/chat-server/internal/service_error"
)

func TestImplementation_SendMessage(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService
	type messageServiceMockFunc func(mc *minimock.Controller) service.MessageService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatIDSuccess    = gofakeit.Int64()
		userIDSuccess    = gofakeit.Int64()
		fromSuccess      = gofakeit.Name()
		textSuccess      = gofakeit.BeerAlcohol()
		timestampSuccess = &timestamppb.Timestamp{}
		chatIDFailed     = 0
		userIDFailed     = 0
		fromFailed       = ""
		textFailed       = ""

		serviceErrCase = serviceErrors.ErrRequireParam
		serviceErr     = fmt.Errorf("service error")

		reqSuccess = &desc.SendMessageRequest{
			ChatId:    chatIDSuccess,
			From:      fromSuccess,
			UserId:    userIDSuccess,
			Text:      textSuccess,
			Timestamp: timestampSuccess,
		}

		reqFailedCase1 = &desc.SendMessageRequest{
			ChatId:    int64(chatIDFailed),
			From:      fromSuccess,
			UserId:    userIDSuccess,
			Text:      textSuccess,
			Timestamp: timestampSuccess,
		}

		reqFailedCase2 = &desc.SendMessageRequest{
			ChatId:    chatIDSuccess,
			From:      fromSuccess,
			UserId:    int64(userIDFailed),
			Text:      textSuccess,
			Timestamp: timestampSuccess,
		}

		reqFailedCase3 = &desc.SendMessageRequest{
			ChatId:    chatIDSuccess,
			From:      fromFailed,
			UserId:    userIDSuccess,
			Text:      textSuccess,
			Timestamp: timestampSuccess,
		}

		reqFailedCase4 = &desc.SendMessageRequest{
			ChatId:    chatIDSuccess,
			From:      fromSuccess,
			UserId:    userIDSuccess,
			Text:      textFailed,
			Timestamp: timestampSuccess,
		}

		message = &model.Message{
			ChatID:    chatIDSuccess,
			UserID:    userIDSuccess,
			From:      fromSuccess,
			Text:      textSuccess,
			Timestamp: timestampSuccess.AsTime(),
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		chatServiceMock    chatServiceMockFunc
		messageServiceMock messageServiceMockFunc
	}{
		{
			name: "success send message",
			args: args{
				ctx: ctx,
				req: reqSuccess,
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(&emptypb.Empty{}, nil)
				return mock
			},
		},
		{
			name: "failed send message",
			args: args{
				ctx: ctx,
				req: reqSuccess,
			},
			want: &emptypb.Empty{},
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(&emptypb.Empty{}, serviceErr)
				return mock
			},
		},
		{
			name: "failed send message case 1",
			args: args{
				ctx: ctx,
				req: reqFailedCase1,
			},
			want: &emptypb.Empty{},
			err:  serviceErrCase,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed send message case 2",
			args: args{
				ctx: ctx,
				req: reqFailedCase2,
			},
			want: &emptypb.Empty{},
			err:  serviceErrCase,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed send message",
			args: args{
				ctx: ctx,
				req: reqFailedCase3,
			},
			want: &emptypb.Empty{},
			err:  serviceErrCase,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				return mock
			},
		},
		{
			name: "failed send message",
			args: args{
				ctx: ctx,
				req: reqFailedCase4,
			},
			want: &emptypb.Empty{},
			err:  serviceErrCase,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := serviceMocks.NewMessageServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatServiceMock(mc)
			messageServiceMock := tt.messageServiceMock(mc)

			api := chat.NewChatServer(chatServiceMock, messageServiceMock)

			newID, err := api.SendMessage(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}

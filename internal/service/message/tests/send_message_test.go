package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/platform_common/pkg/db"

	txMocks "github.com/ArturSaga/chat-server/internal/client/db/mocks"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/repository"
	"github.com/ArturSaga/chat-server/internal/service/message"
	serviceMocks "github.com/ArturSaga/chat-server/internal/service/mocks"
)

func TestImplementation_SendMessage(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx     context.Context
		message *model.Message
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		from      = gofakeit.BeerName()
		userID    = gofakeit.Int64()
		chatID    = gofakeit.Int64()
		text      = gofakeit.CelebrityActor()
		timestamp = time.Now()

		//id           = gofakeit.Int64()
		modelMessage = &model.Message{
			From:      from,
			ChatID:    chatID,
			UserID:    userID,
			Text:      text,
			Timestamp: timestamp,
		}

		serviceErr = fmt.Errorf("service error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name              string
		args              args
		want              *emptypb.Empty
		err               error
		chatRepository    chatRepositoryMockFunc
		messageRepository messageRepositoryMockFunc
		txManager         transactionMockFunc
	}{
		{
			name: "success send",
			args: args{
				ctx:     ctx,
				message: modelMessage,
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageRepository: func(mc *minimock.Controller) repository.MessageRepository {
				mock := serviceMocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, modelMessage).Return(&emptypb.Empty{}, nil)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "failed send",
			args: args{
				ctx:     ctx,
				message: modelMessage,
			},
			want: &emptypb.Empty{},
			err:  serviceErr,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
			messageRepository: func(mc *minimock.Controller) repository.MessageRepository {
				mock := serviceMocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, modelMessage).Return(&emptypb.Empty{}, serviceErr)
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			messageRepositoryMock := tt.messageRepository(mc)
			txManagerMock := tt.txManager(mc)
			messageService := message.NewMessageService(messageRepositoryMock, txManagerMock)

			res, err := messageService.SendMessage(tt.args.ctx, tt.args.message)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}

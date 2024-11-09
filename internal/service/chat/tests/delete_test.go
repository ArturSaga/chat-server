package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/client/db"
	txMocks "github.com/ArturSaga/chat-server/internal/client/db/mocks"
	"github.com/ArturSaga/chat-server/internal/repository"
	"github.com/ArturSaga/chat-server/internal/service/chat"
	serviceMocks "github.com/ArturSaga/chat-server/internal/service/mocks"
)

func TestImplementation_DeleteChat(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		// Пример ошибки
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
			name: "success delete",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc) // Используем правильный мок репозитория
				mock.DeleteChatMock.Expect(ctx, id).Return(&emptypb.Empty{}, nil)
				return mock
			},
			messageRepository: func(mc *minimock.Controller) repository.MessageRepository {
				mock := serviceMocks.NewMessageServiceMock(mc) // Мок для репозитория сообщений
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc) // Мок для транзакционного менеджера
				return mock
			},
		},
		{
			name: "failed delete",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: &emptypb.Empty{},
			err:  serviceErr,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc)                 // Используем правильный мок репозитория
				mock.DeleteChatMock.Expect(ctx, id).Return(nil, serviceErr) // Ошибка при удалении
				return mock
			},
			messageRepository: func(mc *minimock.Controller) repository.MessageRepository {
				mock := serviceMocks.NewMessageServiceMock(mc) // Мок для репозитория сообщений
				return mock
			},
			txManager: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc) // Мок для транзакционного менеджера
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatRepositoryMock := tt.chatRepository(mc)
			//messageRepositoryMock := tt.messageRepository(mc)
			txManagerMock := tt.txManager(mc)

			// Создаем сервис с моками
			chatService := chat.NewChatService(chatRepositoryMock, txManagerMock)

			// Выполняем удаление
			res, err := chatService.DeleteChat(tt.args.ctx, tt.args.id)

			// Проверяем результат
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}

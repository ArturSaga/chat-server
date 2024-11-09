package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ArturSaga/chat-server/internal/client/db"
	txMocks "github.com/ArturSaga/chat-server/internal/client/db/mocks"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/repository"
	"github.com/ArturSaga/chat-server/internal/service/chat"
	serviceMocks "github.com/ArturSaga/chat-server/internal/service/mocks"
)

func TestImplementation_CreateChat(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository
	type transactionMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		chat *model.Chat
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userIDsSuccess   = []int64{1, 2}
		userNamesSuccess = []string{"Kolya", "Petya"}
		chatNameSuccess  = "Test chat"

		id        = gofakeit.Int64()
		chatModel = &model.Chat{
			UserIDs:   userIDsSuccess,
			UserNames: userNamesSuccess,
			ChatName:  chatNameSuccess,
		}

		// Пример ошибки
		serviceErr = fmt.Errorf("service error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name              string
		args              args
		want              int64
		err               error
		chatRepository    chatRepositoryMockFunc
		messageRepository messageRepositoryMockFunc
		txManager         transactionMockFunc
	}{
		{
			name: "success create",
			args: args{
				ctx:  ctx,
				chat: chatModel,
			},
			want: id,
			err:  nil,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc) // Используем правильный мок репозитория
				mock.CreateChatMock.Expect(ctx, chatModel).Return(id, nil)
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
				ctx:  ctx,
				chat: chatModel,
			},
			want: 0,
			err:  serviceErr,
			chatRepository: func(mc *minimock.Controller) repository.ChatRepository {
				mock := serviceMocks.NewChatServiceMock(mc)                      // Используем правильный мок репозитория
				mock.CreateChatMock.Expect(ctx, chatModel).Return(0, serviceErr) // Ошибка при удалении
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
			res, err := chatService.CreateChat(tt.args.ctx, tt.args.chat)

			// Проверяем результат
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}

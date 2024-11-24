package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/api/chat"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/service"
	serviceMocks "github.com/ArturSaga/chat-server/internal/service/mocks"
	serviceErrors "github.com/ArturSaga/chat-server/internal/service_error"
)

func TestImplementation_CreateChat(t *testing.T) {
	t.Parallel()

	// Инициализация контекста и вспомогательных данных
	ctx := context.Background()
	mc := minimock.NewController(t)

	// Генерация тестовых данных с помощью gofakeit
	id := gofakeit.Int64()

	// Данные для успешных и неудачных случаев
	userIDsSuccess := []int64{1, 2}
	userNamesSuccess := []string{"Kolya", "Petya"}
	chatNameSuccess := "Test chat"

	userIDsFailedCase1 := []int64{1, 2}
	userNamesFailedCase1 := []string{"Petya"}
	chatNameFailedCase2 := ""
	serviceErrCase1 := serviceErrors.ErrUserIDsNotMatchUserNames
	serviceErrCase2 := serviceErrors.ErrRequireParam
	serviceErr := fmt.Errorf("service error")

	// Создание запросов для тестов
	reqSuccess := &desc.CreateChatRequest{
		UserIds:   userIDsSuccess,
		Usernames: userNamesSuccess,
		ChatName:  chatNameSuccess,
	}

	reqFailedCase1 := &desc.CreateChatRequest{
		UserIds:   userIDsFailedCase1,
		Usernames: userNamesFailedCase1,
		ChatName:  chatNameSuccess,
	}

	reqFailedCase2 := &desc.CreateChatRequest{
		UserIds:   userIDsSuccess,
		Usernames: userNamesSuccess,
		ChatName:  chatNameFailedCase2,
	}

	reqFailedCase3 := &desc.CreateChatRequest{
		UserIds:   nil,
		Usernames: userNamesSuccess,
		ChatName:  chatNameSuccess,
	}

	reqFailedCase4 := &desc.CreateChatRequest{
		UserIds:   userIDsSuccess,
		Usernames: nil,
		ChatName:  chatNameSuccess,
	}

	chatModelSuccess := &model.Chat{
		UserIDs:   userIDsSuccess,
		UserNames: userNamesSuccess,
		ChatName:  chatNameSuccess,
	}

	res := &desc.CreateChatResponse{
		Id: id,
	}

	// Очистка контроллера после завершения теста
	t.Cleanup(mc.Finish)

	// Группировка тестов
	tests := []struct {
		name string
		args struct {
			ctx context.Context
			req *desc.CreateChatRequest
		}
		want               *desc.CreateChatResponse
		err                error
		chatServiceMock    func(mc *minimock.Controller) service.ChatService
		messageServiceMock func(mc *minimock.Controller) service.MessageService
	}{
		// Успешный тест
		{
			name: "success create",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqSuccess,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, chatModelSuccess).Return(id, nil)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
		{
			name: "failed create",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqSuccess,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, chatModelSuccess).Return(0, serviceErr)
				return mock
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
		{
			name: "failed validate case 1",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqFailedCase1,
			},
			want: nil,
			err:  serviceErrCase1,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				// В этом случае сервис не должен быть вызван
				return serviceMocks.NewChatServiceMock(mc)
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
		// Тест с ошибкой: пустое имя чата
		{
			name: "failed validate case 2",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqFailedCase2,
			},
			want: nil,
			err:  serviceErrCase2,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
		// Тест с ошибкой: отсутствует UserIds
		{
			name: "failed validate case 3",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqFailedCase3,
			},
			want: nil,
			err:  serviceErrCase2,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
		// Тест с ошибкой: отсутствуют Usernames
		{
			name: "failed validate case 4",
			args: struct {
				ctx context.Context
				req *desc.CreateChatRequest
			}{
				ctx: ctx,
				req: reqFailedCase4,
			},
			want: nil,
			err:  serviceErrCase2,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				return serviceMocks.NewMessageServiceMock(mc)
			},
		},
	}

	// Выполнение тестов
	for _, tt := range tests {
		tt := tt // Создаем копию для каждой итерации (для параллельных тестов)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создаем моки для ChatService и MessageService
			chatServiceMock := tt.chatServiceMock(mc)
			messageServiceMock := tt.messageServiceMock(mc)

			// Создаем объект API
			api := chat.NewChatServer(chatServiceMock, messageServiceMock)

			// Вызов функции CreateChat
			newID, err := api.CreateChat(tt.args.ctx, tt.args.req)

			// Проверка результата
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}

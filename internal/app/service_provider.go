package app

import (
	"context"
	"log"

	"github.com/ArturSaga/chat-server/internal/api/chat"
	"github.com/ArturSaga/chat-server/internal/client/db"
	"github.com/ArturSaga/chat-server/internal/client/db/pg"
	"github.com/ArturSaga/chat-server/internal/client/db/transaction"
	"github.com/ArturSaga/chat-server/internal/closer"
	"github.com/ArturSaga/chat-server/internal/config"
	"github.com/ArturSaga/chat-server/internal/repository"
	chatRepository "github.com/ArturSaga/chat-server/internal/repository/chat"
	messageRepository "github.com/ArturSaga/chat-server/internal/repository/message"
	"github.com/ArturSaga/chat-server/internal/service"
	chatService "github.com/ArturSaga/chat-server/internal/service/chat"
	"github.com/ArturSaga/chat-server/internal/service/message"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient          db.Client
	txManager         db.TxManager
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository

	chatService    service.ChatService
	messageService service.MessageService

	implementation *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s serviceProvider) PgConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewChatRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = messageRepository.NewMessageRepository(s.DBClient(ctx))
	}

	return s.messageRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewChatService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) MessageService(ctx context.Context) service.MessageService {
	if s.messageService == nil {
		s.messageService = message.NewMessageService(
			s.MessageRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.messageService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.implementation == nil {
		s.implementation = chat.NewImplementation(s.ChatService(ctx), s.MessageService(ctx))
	}

	return s.implementation
}

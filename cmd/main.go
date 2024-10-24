package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/ArturSaga/chat-server/api/grpc/pkg/chat_v1"
	"github.com/ArturSaga/chat-server/internal/config"
	"github.com/ArturSaga/chat-server/internal/config/env"
)

type server struct {
	desc.UnimplementedChatApiServer
	pool *pgxpool.Pool
}

const timeout = 15

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if _, ok := ctx.Deadline(); !ok {
		c, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		ctx = c
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatApiServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateChat - публичный метод, который создает новый чат.
func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	builderChatInsert := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_name", "created_at").
		Values(req.ChatName, time.Now()).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, err
	}

	var chatID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		fmt.Printf("failed to insert chat: %v", err)
		return nil, err
	}

	builderChatUsersInsert := sq.Insert("chat_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "user_name")

	for i := 0; i < len(req.UserIds); i++ {
		builderChatUsersInsert = builderChatUsersInsert.Values(chatID, req.UserIds[i], req.Usernames[i])
	}
	query, args, err = builderChatUsersInsert.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("failed to insert chat_users: %v", err)
		return nil, err
	}

	return &desc.CreateChatResponse{Id: chatID}, nil
}

// SendMessage - публичный метод, который отправляет сообщение в чат.
func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	builderMessageInsert := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id", "user_name", "text", "timestamp").
		Values(req.ChatId, req.UserId, req.From, req.Text, time.Now())

	query, args, err := builderMessageInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, nil
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("failed to send message: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteChat - публичный метод, который удаляет чат.
func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	builderDeleteChatUsers := sq.Delete("chat_users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"chat_id": req.Id})
	query, args, err := builderDeleteChatUsers.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, nil
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("failed to delete chat_users: %v", err)
		return nil, err
	}

	builderDeleteChat := sq.Delete("chats").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})
	query, args, err = builderDeleteChat.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, nil
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("failed to delete chats: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/client/db"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/repository"
)

const chatTableName = "chats"
const userChatTableName = "chat_users"

const idChatColumn = "id"
const nameChatColumn = "chat_name"
const createdAtChatColumn = "created_at"

const chatIdColumn = "chat_id"
const userIdColumn = "user_id"
const userNameColumn = "user_name"

type repo struct {
	db db.Client
}

// NewChatRepository - публичный метод, создащий сущность репозитория, для работы с данными сущности в бд
func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	// Делаем запрос на вставку записи в таблицу user
	builderChatInsert := sq.Insert(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameChatColumn).
		Values(chat.ChatName).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		fmt.Printf("failed to insert chat: %v", err)
		return 0, err
	}

	builderUserChatInsert := sq.Insert(userChatTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIdColumn, userIdColumn, userNameColumn)

	for i := range chat.UserIDs {
		builderUserChatInsert = builderUserChatInsert.Values(chatID, chat.UserIDs[i], chat.UserNames[i])
	}

	query, args, err = builderUserChatInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return 0, err
	}

	q = db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		fmt.Printf("failed to insert chat: %v", err)
		return 0, err
	}

	log.Printf("inserted chat with id: %d", chatID)

	return chatID, nil
}

func (r *repo) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	builderUserDelete := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idChatColumn: id})

	query, args, err := builderUserDelete.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return &emptypb.Empty{}, err
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		fmt.Printf("failed to delete chat: %v", err)
	}

	return &emptypb.Empty{}, nil
}

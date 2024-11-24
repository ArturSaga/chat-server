package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/platform_common/pkg/db"

	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/repository"
)

const chatTableName = "chats"
const userChatTableName = "chat_users"

const idColumnChat = "id"
const nameColumnChat = "chat_name"

const chatIDColumnUserChat = "chat_id"
const userIDColumnUserChat = "user_id"
const userNameColumnUserChat = "user_name"

type repo struct {
	db db.Client
}

// NewChatRepository - публичный метод, создащий сущность репозитория, для работы с данными сущности в бд
func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// CreateChat - публичный метод, для создания чата в слое репозитория
func (r *repo) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	// Делаем запрос на вставку записи в таблицу chats
	builderChatInsert := sq.Insert(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumnChat).
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
		Columns(chatIDColumnUserChat, userIDColumnUserChat, userNameColumnUserChat)

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

// DeleteChat - публичный метод, для удаления чата в слое репозитория
func (r *repo) DeleteChat(ctx context.Context, id int64) (*emptypb.Empty, error) {
	builderUserDelete := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumnChat: id})

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

package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/chat-server/internal/client/db"
	"github.com/ArturSaga/chat-server/internal/model"
	"github.com/ArturSaga/chat-server/internal/repository"
)

const tableName = "messages"

const userNameColumn = "user_name"
const chatIDColumn = "chat_id"
const userIDColumn = "user_id"
const textColumn = "text"
const timestampColumn = "timestamp"

type repo struct {
	db db.Client
}

// NewMessageRepository - публичный метод, создащий сущность репозитория, для работы с данными сущности в бд
func NewMessageRepository(db db.Client) repository.MessageRepository {
	return &repo{db: db}
}

// SendMessage - публичный метод, для отправки сообщения в слое репозитория
func (r *repo) SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error) {
	fmt.Println(message.Text)
	// Делаем запрос на вставку записи в таблицу user
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, userIDColumn, userNameColumn, textColumn, timestampColumn).
		Values(message.ChatID, message.UserID, message.From, message.Text, message.Timestamp)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return &emptypb.Empty{}, err
	}

	q := db.Query{
		Name:     "message_repository.SendMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		fmt.Printf("failed to insert message: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

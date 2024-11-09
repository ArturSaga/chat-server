package model

import "time"

// Message - сущность сообщения, для работы с сервисным слоем
type Message struct {
	From      *string
	ChatID    *int64
	UserID    *int64
	Text      *string
	Timestamp *time.Time
}

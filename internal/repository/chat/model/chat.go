package model

// Chat - сущность чата, для работы с сервисным слоем
type Chat struct {
	UserIDs   []int64
	UserNames []string
	ChatName  string
}

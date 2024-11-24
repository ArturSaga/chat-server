package service_error

import "errors"

var (
	// ErrUserIDsNotMatchUserNames - ошибка не соответствия паролей
	ErrUserIDsNotMatchUserNames = errors.New("count userIDs not match to count usernames")
	// ErrRequireParam - ошибка получения данных пользователя
	ErrRequireParam = errors.New("one of the parameters is nil")
	// ErrUpdateUser - ошибка при обновлении данных пользователя
)

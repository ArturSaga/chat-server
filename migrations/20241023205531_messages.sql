-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
  id SERIAL PRIMARY KEY,
  chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
  user_id BIGINT,  -- идентификатор пользователя, который может ссылаться на другой сервис
  user_name VARCHAR(255),
  text TEXT,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS messages;
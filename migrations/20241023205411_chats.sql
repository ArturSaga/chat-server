-- +goose Up
CREATE TABLE IF NOT EXISTS chats (
   id SERIAL PRIMARY KEY,
   chat_name VARCHAR(255),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS chats;
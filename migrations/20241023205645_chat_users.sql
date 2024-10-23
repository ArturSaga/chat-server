-- +goose Up
CREATE TABLE IF NOT EXISTS chat_users (
    chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT,
    user_name VARCHAR(255),
    PRIMARY KEY (chat_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS chat_users;

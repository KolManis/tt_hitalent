-- +goose Up
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);

-- +goose Down
DROP TABLE messages;
-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE chats;
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd

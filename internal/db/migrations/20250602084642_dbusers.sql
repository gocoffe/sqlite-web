-- +goose Up
-- +goose StatementBegin
CREATE TABLE dbusers
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    login    TEXT NOT NULL,
    password TEXT NOT NULL,

    UNIQUE (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbusers;
-- +goose StatementEnd

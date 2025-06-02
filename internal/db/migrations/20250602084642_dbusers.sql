-- +goose Up
-- +goose StatementBegin
CREATE TABLE dbusers
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    login    TEXT NOT NULL,
    password TEXT NOT NULL,

    UNIQUE (login)
);

-- TODO: develop other method changing default password (also should be hashed)
INSERT INTO dbusers (login, password)
VALUES ('admin', 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbusers;
-- +goose StatementEnd

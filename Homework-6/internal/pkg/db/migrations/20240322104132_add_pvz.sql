-- +goose Up
-- +goose StatementBegin
CREATE TABLE pvz(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz;
-- +goose StatementEnd

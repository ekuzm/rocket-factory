-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    uuid UUID NOT NULL,
    user_uuid UUID NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    transaction_uuid UUID NULL,
    status INTEGER NOT NULL CHECK(status BETWEEN 0 AND 2),
    payment_method INTEGER CHECK(status BETWEEN 0 AND 4),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    PRIMARY KEY(uuid)
);
-- +goose StatementEnd
DROP TABLE IF EXISTS orders;
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

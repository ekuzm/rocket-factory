-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    uuid UUID NOT NULL,
    user_uuid UUID NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    transaction_uuid UUID NULL,
    status INTEGER NOT NULL CHECK(status BETWEEN 0 AND 2),
    payment_method INTEGER CHECK(payment_method BETWEEN 0 AND 4),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    PRIMARY KEY(uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd

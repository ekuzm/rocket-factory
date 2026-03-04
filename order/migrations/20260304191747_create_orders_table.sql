-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids UUID ARRAY,
    total_price DOUBLE PRECISION NOT NULL,
    transaction_uuid UUID,
    status order_status NOT NULL,
    payment_method payment_method,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    UNIQUE (user_uuid, transaction_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd

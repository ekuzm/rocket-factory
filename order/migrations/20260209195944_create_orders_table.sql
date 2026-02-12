-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    transaction_uuid UUID NOT NULL,
    payment_method VARCHAR(25) NOT NULL DEFAULT 'unknown',
    status VARCHAR(25) NOT NULL DEFAULT 'pending_payment',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
    CREATE TYPE payment_method AS ENUM('unknown', 'card', 'spb', 'credit_card', 'investor_money');
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS payment_method;
-- +goose StatementEnd

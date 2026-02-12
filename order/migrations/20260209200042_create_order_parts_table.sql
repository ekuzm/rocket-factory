-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_parts (
    order_uuid UUID,
    part_uuid UUID,
    PRIMARY KEY(order_uuid, part_uuid),
    FOREIGN KEY(order_uuid) REFERENCES orders(uuid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_parts;
-- +goose StatementEnd

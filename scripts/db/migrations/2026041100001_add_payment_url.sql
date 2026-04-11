-- +migrate Up
ALTER TABLE repair_orders ADD COLUMN payment_url TEXT;

-- +migrate Down
ALTER TABLE repair_orders DROP COLUMN payment_url;

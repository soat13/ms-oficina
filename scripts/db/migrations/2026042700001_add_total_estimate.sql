-- +migrate Up
ALTER TABLE repair_orders ADD COLUMN total_estimate BIGINT;

-- +migrate Down
ALTER TABLE repair_orders DROP COLUMN total_estimate;

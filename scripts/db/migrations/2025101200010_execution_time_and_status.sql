-- +migrate Up
ALTER TABLE repair_orders ADD COLUMN execution_time_minutes BIGINT;

UPDATE repair_orders 
SET status = 'finished' 
WHERE id = '3b5d2e78-9f0a-47b5-9a1c-6e8f2b3d4c5e' AND status = 'completed';

-- +migrate Down
ALTER TABLE repair_orders DROP COLUMN execution_time_minutes;

UPDATE repair_orders 
SET status = 'completed' 
WHERE id = '3b5d2e78-9f0a-47b5-9a1c-6e8f2b3d4c5e' AND status = 'finished';

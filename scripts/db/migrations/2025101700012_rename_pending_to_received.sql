-- +migrate Up
UPDATE repair_orders 
SET status = 'received' 
WHERE status = 'pending';

-- +migrate Down
UPDATE repair_orders 
SET status = 'pending' 
WHERE status = 'received';


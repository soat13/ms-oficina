-- +migrate Up
ALTER TABLE estimates RENAME COLUMN repair_id TO repair_order_id;
ALTER INDEX idx_estimates_repair_id RENAME TO idx_estimates_repair_order_id;

UPDATE repair_orders 
SET status = 'awaiting_approval' 
WHERE id = '2c4a1f56-76b3-4d88-8c2f-5b7e9d1f0c3a' AND status = 'in_progress';

UPDATE estimates 
SET status = 'rejected' 
WHERE id = '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a' AND status = 'draft';

UPDATE estimates 
SET status = 'awaiting_approval' 
WHERE id = '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b' AND status = 'sent';

-- +migrate Down
ALTER TABLE estimates RENAME COLUMN repair_order_id TO repair_id;
ALTER INDEX idx_estimates_repair_order_id RENAME TO idx_estimates_repair_id;

UPDATE repair_orders 
SET status = 'in_progress' 
WHERE id = '2c4a1f56-76b3-4d88-8c2f-5b7e9d1f0c3a' AND status = 'awaiting_approval';

UPDATE estimates 
SET status = 'draft' 
WHERE id = '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a' AND status = 'rejected';

UPDATE estimates 
SET status = 'sent' 
WHERE id = '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b' AND status = 'awaiting_approval';


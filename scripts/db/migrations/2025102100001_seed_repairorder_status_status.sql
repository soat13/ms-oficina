-- +migrate Up
CREATE TYPE repairorder_status AS ENUM ('received', 'in_diagnostics', 'awaiting_approval', 'approved', 'in_execution', 'finished', 'released');

-- +migrate Down
DROP TYPE IF EXISTS repairorder_status;
-- +migrate Up
CREATE INDEX idx_vehicles_plate ON vehicles (plate);

-- +migrate Down
DROP INDEX IF EXISTS idx_vehicles_plate;

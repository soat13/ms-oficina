-- +migrate Up
UPDATE vehicles SET plate = UPPER(plate);
UPDATE products SET name = LOWER(name);
UPDATE services SET name = LOWER(name);


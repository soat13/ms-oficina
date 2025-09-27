-- +migrate Up
INSERT INTO services (id, name, price)
VALUES ('1f094799-f210-63a0-a800-e67186b3a9f6', 'Troca de Óleo', 15000),
       ('1f09479a-5f13-6f10-a20a-bde566d4cfd9', 'Alinhamento e Balanceamento', 12000)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, name, price, stock)
VALUES ('1f09479b-6074-69d0-a0d0-9a7c61ccc8bb', 'Filtro de Óleo', 4500, 50),
       ('1f09479b-bc7c-62f0-a9bc-9390bea2885f', 'Pneu 205/55 R16', 350000, 20)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO customers
(id, "name", cpf_cnpj, created_at, updated_at)
VALUES('321e4567-e89b-12d3-a456-426614174000','Ragnar Lodbrok','09807604321', now(), now());
    ON CONFLICT (id) DO NOTHING

INSERT INTO vehicles
(id, customer_id, plate, brand, model, "year", created_at, updated_at)
VALUES
('123e4567-e89b-12d3-a456-426614174000', '321e4567-e89b-12d3-a456-426614174000', 'abc-1234', 'honda', 'corola', 2020, now(), now());
    ON CONFLICT (id) DO NOTHING;

-- +migrate Down
DELETE FROM products
WHERE id IN ('1f09479b-6074-69d0-a0d0-9a7c61ccc8bb','1f09479b-bc7c-62f0-a9bc-9390bea2885f');

DELETE FROM services
WHERE id IN ('1f094799-f210-63a0-a800-e67186b3a9f6','1f09479a-5f13-6f10-a20a-bde566d4cfd9');

DELETE FROM vehicles WHERE id IN ('123e4567-e89b-12d3-a456-426614174000')

DELETE FROM customers WHERE id IN ('321e4567-e89b-12d3-a456-426614174000')

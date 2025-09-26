-- +migrate Up
CREATE TABLE customers
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(255)       NOT NULL,
    cpf_cnpj   VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP          NOT NULL DEFAULT NOW()
);

CREATE TABLE vehicles
(
    id          UUID PRIMARY KEY,
    customer_id UUID         NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    plate       VARCHAR(10)  NOT NULL UNIQUE,
    brand       VARCHAR(100) NOT NULL,
    model       VARCHAR(100) NOT NULL,
    year        INT          NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_vehicles_customer_id ON vehicles (customer_id);

CREATE TABLE services
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    price      BIGINT       NOT NULL,
    currency   VARCHAR(3)   NOT NULL DEFAULT 'BRL',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_services_name ON services (name);

CREATE TABLE products
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    price      BIGINT       NOT NULL,
    stock      INT          NOT NULL DEFAULT 0,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_products_name ON products (name);

CREATE TABLE repair_orders
(
    id          UUID PRIMARY KEY,
    customer_id UUID        NOT NULL REFERENCES customers (id),
    vehicle_id  UUID        NOT NULL REFERENCES vehicles (id),
    status      VARCHAR(30) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_repair_orders_customer_id ON repair_orders (customer_id);
CREATE INDEX idx_repair_orders_vehicle_id ON repair_orders (vehicle_id);
CREATE INDEX idx_repair_orders_status ON repair_orders (status);
CREATE INDEX idx_repair_orders_created_at ON repair_orders (created_at);

CREATE TABLE estimates
(
    id         UUID PRIMARY KEY,
    repair_id  UUID        NOT NULL REFERENCES repair_orders (id) ON DELETE CASCADE,
    status     VARCHAR(30) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_estimates_repair_id ON estimates (repair_id);
CREATE INDEX idx_estimates_status ON estimates (status);
CREATE INDEX idx_estimates_created_at ON estimates (created_at);

CREATE TABLE estimate_items
(
    id          UUID PRIMARY KEY,
    estimate_id UUID         NOT NULL REFERENCES estimates (id) ON DELETE CASCADE,
    item_id     UUID         NOT NULL,
    item_name   VARCHAR(255) NOT NULL,
    item_type   VARCHAR(20)  NOT NULL,
    price       BIGINT       NOT NULL,
    quantity    INT          NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_estimate_items_estimate_id ON estimate_items (estimate_id);
CREATE INDEX idx_estimate_items_item ON estimate_items (item_id, item_type);

-- +migrate Down
DROP TABLE IF EXISTS estimate_items;
DROP TABLE IF EXISTS estimates;
DROP TABLE IF EXISTS repair_orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS customers;
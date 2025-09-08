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

CREATE TABLE services
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    price_cents BIGINT       NOT NULL,
    currency    VARCHAR(3)   NOT NULL DEFAULT 'BRL',
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE products
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    price_cents BIGINT       NOT NULL,
    stock       INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE repair_orders
(
    id          UUID PRIMARY KEY,
    customer_id UUID        NOT NULL REFERENCES customers (id),
    vehicle_id  UUID        NOT NULL REFERENCES vehicles (id),
    status      VARCHAR(30) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE estimates
(
    id         UUID PRIMARY KEY,
    repair_id  UUID        NOT NULL REFERENCES repair_orders (id) ON DELETE CASCADE,
    status     VARCHAR(30) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE estimate_items
(
    id          UUID PRIMARY KEY,
    estimate_id UUID         NOT NULL REFERENCES estimates (id) ON DELETE CASCADE,
    item_id     UUID         NOT NULL,
    item_name   VARCHAR(255) NOT NULL,
    item_type   VARCHAR(20)  NOT NULL,-- todo: ENUM('service', 'product')
    price_cents BIGINT       NOT NULL,
    quantity    INT          NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

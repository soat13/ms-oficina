-- +migrate Up
INSERT INTO public.services (
    id,
    name,
    price,
    currency
) VALUES
    (
        '020b87d3-5c7d-4c66-9d53-4bc7d2b1ad9a',
        'Troca de Óleo',
        15000,
        'BRL'
    ),
    (
        '56ef9a12-2c82-4a37-8fbd-b6e1bd1e2efc',
        'Alinhamento e Balanceamento',
        20000,
        'BRL'
    ),
    (
        '9ac4ef89-7841-4b10-b9bb-9b6200fd1c72',
        'Revisão Completa',
        45000,
        'BRL'
    );

-- +migrate Down
DELETE FROM public.services
WHERE id IN (
    '020b87d3-5c7d-4c66-9d53-4bc7d2b1ad9a',
    '56ef9a12-2c82-4a37-8fbd-b6e1bd1e2efc',
    '9ac4ef89-7841-4b10-b9bb-9b6200fd1c72'
);

-- +migrate Up
INSERT INTO public.estimate_items (
    id,
    estimate_id,
    item_id,
    item_name,
    item_type,
    price,
    quantity
) VALUES
    (
        '701c6b24-b5e9-4c9d-c4e5-5b7c9d0e1f2a',
        '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a',
        '020b87d3-5c7d-4c66-9d53-4bc7d2b1ad9a',
        'Troca de Óleo',
        'service',
        15000,
        1
    ),
    (
        '812d7c35-c6fa-4d1e-d5f6-6c8d0e1f2a3b',
        '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a',
        'f2b1aeb3-0d0a-465a-b91e-6122e3aa0baf',
        'Filtro de Óleo',
        'product',
        4500,
        1
    ),
    (
        '923e8d46-d70b-4e2f-e607-7d9e0f1a2b3c',
        '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b',
        '56ef9a12-2c82-4a37-8fbd-b6e1bd1e2efc',
        'Alinhamento e Balanceamento',
        'service',
        20000,
        1
    ),
    (
        'a34f9e57-e81c-4f30-f718-8e0f1a2b3c4d',
        '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b',
        '67c7c2f2-55ac-4fbf-8caa-7c7d50e3f0b2',
        'Pastilha de Freio',
        'product',
        8000,
        1
    ),
    (
        'b450af68-f92d-5031-0829-9f102a3b4c5d',
        '6f9b5a13-a4d8-4b8c-b3d4-4a6b8f9d0e1c',
        '9ac4ef89-7841-4b10-b9bb-9b6200fd1c72',
        'Revisão Completa',
        'service',
        45000,
        1
    ),
    (
        'c561b079-0a3e-5142-193a-af213b4c5d6e',
        '6f9b5a13-a4d8-4b8c-b3d4-4a6b8f9d0e1c',
        'a4e1c9bd-41a8-4ac7-a0e0-5a2d4fce6d32',
        'Velas de Ignição',
        'product',
        6000,
        4
    );

-- +migrate Down
DELETE FROM public.estimate_items
WHERE id IN (
    '701c6b24-b5e9-4c9d-c4e5-5b7c9d0e1f2a',
    '812d7c35-c6fa-4d1e-d5f6-6c8d0e1f2a3b',
    '923e8d46-d70b-4e2f-e607-7d9e0f1a2b3c',
    'a34f9e57-e81c-4f30-f718-8e0f1a2b3c4d',
    'b450af68-f92d-5031-0829-9f102a3b4c5d',
    'c561b079-0a3e-5142-193a-af213b4c5d6e'
);

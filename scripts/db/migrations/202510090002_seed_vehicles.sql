-- +migrate Up
INSERT INTO public.vehicles (
    id,
    customer_id,
    plate,
    brand,
    model,
    year
) VALUES
    (
        'a3b834e2-0d60-4c7d-9fd0-2966b5aa1dba',
        '8f61c381-663b-4bf1-93c8-67a0a28cceaf',
        'ABC1D23',
        'Ford',
        'Fiesta',
        2018
    ),
    (
        '893ab8bb-4f84-487c-9f3f-1f63d5e4b85e',
        '2e72f0ed-2b32-4a07-a6f8-b21c8df6e2d1',
        'XYZ4E56',
        'Chevrolet',
        'Onix',
        2020
    ),
    (
        '5f7ae479-4040-43d1-8a8d-67c4bc51fd12',
        'dd1e143a-5e86-4de3-9ed0-d9fe380395c2',
        'JKL7M89',
        'Toyota',
        'Corolla',
        2021
    );

-- +migrate Down
DELETE FROM public.vehicles
WHERE id IN (
    'a3b834e2-0d60-4c7d-9fd0-2966b5aa1dba',
    '893ab8bb-4f84-487c-9f3f-1f63d5e4b85e',
    '5f7ae479-4040-43d1-8a8d-67c4bc51fd12'
);

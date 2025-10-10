-- +migrate Up
INSERT INTO public.repair_orders (
    id,
    customer_id,
    vehicle_id,
    status
) VALUES
    (
        '1f3e1c89-3ad4-4ca9-87a9-0a6f1dcb1e2f',
        '8f61c381-663b-4bf1-93c8-67a0a28cceaf',
        'a3b834e2-0d60-4c7d-9fd0-2966b5aa1dba',
        'pending'
    ),
    (
        '2c4a1f56-76b3-4d88-8c2f-5b7e9d1f0c3a',
        '2e72f0ed-2b32-4a07-a6f8-b21c8df6e2d1',
        '893ab8bb-4f84-487c-9f3f-1f63d5e4b85e',
        'in_progress'
    ),
    (
        '3b5d2e78-9f0a-47b5-9a1c-6e8f2b3d4c5e',
        'dd1e143a-5e86-4de3-9ed0-d9fe380395c2',
        '5f7ae479-4040-43d1-8a8d-67c4bc51fd12',
        'completed'
    );

-- +migrate Down
DELETE FROM public.repair_orders
WHERE id IN (
    '1f3e1c89-3ad4-4ca9-87a9-0a6f1dcb1e2f',
    '2c4a1f56-76b3-4d88-8c2f-5b7e9d1f0c3a',
    '3b5d2e78-9f0a-47b5-9a1c-6e8f2b3d4c5e'
);

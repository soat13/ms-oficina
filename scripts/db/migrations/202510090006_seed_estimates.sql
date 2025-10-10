-- +migrate Up
INSERT INTO public.estimates (
    id,
    repair_id,
    status
) VALUES
    (
        '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a',
        '1f3e1c89-3ad4-4ca9-87a9-0a6f1dcb1e2f',
        'draft'
    ),
    (
        '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b',
        '2c4a1f56-76b3-4d88-8c2f-5b7e9d1f0c3a',
        'sent'
    ),
    (
        '6f9b5a13-a4d8-4b8c-b3d4-4a6b8f9d0e1c',
        '3b5d2e78-9f0a-47b5-9a1c-6e8f2b3d4c5e',
        'approved'
    );

-- +migrate Down
DELETE FROM public.estimates
WHERE id IN (
    '4d7f3e91-82b6-4f6a-9be1-2e4f6d7b8c9a',
    '5e8a4f02-93c7-4a7b-a2c3-3f5a7e8c9d0b',
    '6f9b5a13-a4d8-4b8c-b3d4-4a6b8f9d0e1c'
);

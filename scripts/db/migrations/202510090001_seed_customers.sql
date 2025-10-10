-- +migrate Up
INSERT INTO public.customers (
    id,
    name,
    document,
    document_type,
    email,
    phone_number
) VALUES
    (
        '8f61c381-663b-4bf1-93c8-67a0a28cceaf',
        'João Silva',
        '12345678901',
        'CPF',
        'joao.silva@example.com',
        '11987654321'
    ),
    (
        '2e72f0ed-2b32-4a07-a6f8-b21c8df6e2d1',
        'Maria Oliveira',
        '98765432100',
        'CPF',
        'maria.oliveira@example.com',
        '11999998888'
    ),
    (
        'dd1e143a-5e86-4de3-9ed0-d9fe380395c2',
        'Carlos Souza',
        '45678912355',
        'CPF',
        'carlos.souza@example.com',
        '11988887777'
    );

-- +migrate Down
DELETE FROM public.customers
WHERE id IN (
    '8f61c381-663b-4bf1-93c8-67a0a28cceaf',
    '2e72f0ed-2b32-4a07-a6f8-b21c8df6e2d1',
    'dd1e143a-5e86-4de3-9ed0-d9fe380395c2'
);

-- +migrate Up
INSERT INTO public.products (
    id,
    name,
    price,
    stock
) VALUES
    (
        'f2b1aeb3-0d0a-465a-b91e-6122e3aa0baf',
        'Filtro de Óleo',
        4500,
        25
    ),
    (
        '67c7c2f2-55ac-4fbf-8caa-7c7d50e3f0b2',
        'Pastilha de Freio',
        8000,
        40
    ),
    (
        'a4e1c9bd-41a8-4ac7-a0e0-5a2d4fce6d32',
        'Velas de Ignição',
        6000,
        60
    );

-- +migrate Down
DELETE FROM public.products
WHERE id IN (
    'f2b1aeb3-0d0a-465a-b91e-6122e3aa0baf',
    '67c7c2f2-55ac-4fbf-8caa-7c7d50e3f0b2',
    'a4e1c9bd-41a8-4ac7-a0e0-5a2d4fce6d32'
);

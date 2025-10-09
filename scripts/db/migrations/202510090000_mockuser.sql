-- +migrate Up

INSERT INTO public.users (
    id,
    name,
    document,
    document_type,
    email,
    phone_number,
    password,
    roles
) VALUES (
    'c3d90687-7e00-4c18-8b51-1b4735418bd9',
    'Admin User',
    '71296750043',
    'CPF',
    'admin@example.com',
    '11987654999',
    '$2a$10$leoIh4gK2hfgMLJmpOhJXOndyQB/7z4Xc7sFO1WTrQ2UA6dj6TVBW',
    ARRAY['manager']::public.user_role[]
);

-- +migrate Down
DELETE FROM public.users WHERE id = 'c3d90687-7e00-4c18-8b51-1b4735418bd9';

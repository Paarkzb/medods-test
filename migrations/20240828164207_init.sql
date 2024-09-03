-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.user
(
    id uuid primary key default uuid_generate_v4(),
    name varchar(255) not null,
    email varchar(255) not null unique,
    username varchar(255) not null unique,
    password varchar(255) not null,
    refresh_token varchar(255) null unique,
    refresh_expiration_time timestamp null,
    ip_address varchar(255) null,
    deleted boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.user;
-- +goose StatementEnd

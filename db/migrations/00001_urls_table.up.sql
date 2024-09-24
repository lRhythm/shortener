create extension if not exists "uuid-ossp";

create table if not exists urls
(
    id             uuid         not null default uuid_generate_v4(),
    correlation_id text         not null default '',
    short_url      text         not null,
    original_url   text         not null,
    created_at     timestamp(0) not null default now(),
    updated_at     timestamp(0) not null default now(),
    primary key (id)
);

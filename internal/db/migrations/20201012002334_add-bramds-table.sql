-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists brands
(
    id               bigserial primary key,
    created_at       timestamptz   not null default now(),
    updated_at       timestamptz   not null default now(),
    deleted_at       timestamptz,
    name             varchar(512)  not null default '',
    page_title       varchar(512)  not null default '',
    meta_keywords    varchar[],
    meta_description varchar(2048) not null default '',
    search_keywords  varchar(512)  not null default '',
    custom_url       jsonb default '{}'::jsonb
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists brands;

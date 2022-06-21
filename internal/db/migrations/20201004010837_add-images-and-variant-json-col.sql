-- +goose Up
-- SQL in this section is executed when the migration is applied.
alter table products add images jsonb not null default '{}'::jsonb;
alter table products add variants jsonb not null default '{}'::jsonb;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
alter table products drop images;
alter table products drop variants;

-- name: CreateBrand :one
insert into brands (id,
                    created_at,
                    updated_at,
                    deleted_at,
                    name,
                    page_title,
                    meta_keywords,
                    meta_description,
                    search_keywords,
                    custom_url)

values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)

on conflict (id) do update set updated_at       = now(),
                               name             = $5,
                               page_title       = $6,
                               meta_keywords    = $7,
                               meta_description = $8,
                               search_keywords  = $9,
                               custom_url       = $10
returning *;

-- name: ListBrands :many
select *
from brands;

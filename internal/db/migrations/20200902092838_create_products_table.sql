-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists products
(
    id                         bigserial primary key,
    created_at                 timestamptz not null default now(),
    updated_at                 timestamptz not null default now(),
    deleted_at                 timestamptz,
    name                       varchar(128) not null default '',
    type                       varchar(128) not null default '',
    sku                        varchar(128) not null default '',
    description                varchar(5120) not null default '',
    weight                     double precision not null default 0,
    Depth                      double precision not null default 0,
    Height                     double precision not null default 0,
    Price                      double precision not null default 0,
    CostPrice                  double precision not null default 0,
    RetailPrice                double precision not null default 0,
    SalePrice                  double precision not null default 0,
    MapPrice                   double precision not null default 0,
    tax_class_id               int8 not null default 0,
    product_tax_code           varchar(32) not null default '',
    calculated_price           double precision not null default 0.0,
    categories                 int8[],
    brand_id                   int8 not null default 0,
    option_set_id              int8 not null default 0,
    option_set_display         varchar(32) not null default '',
    inventory_level            int8 not null default 0,
    inventory_warning_level    int8 not null default 0,
    inventory_tracking         varchar(32) not null default '',
    reviews_rating_sum         int8 not null default 0,
    reviews_count              int8 not null default 0,
    total_sold                 int8 not null default 0,
    fixed_cost_shipping_price  double precision not null default 0,
    is_free_shipping           bool not null default false,
    is_visible                 bool not null default false,
    is_featured                bool not null default false,
    warranty                   varchar(256) not null default '',
    bin_picking_number         varchar(128) not null default '',
    layout_file                varchar(128) not null default '',
    upc                        varchar(128) not null default '',
    mpn                        varchar(128) not null default '',
    gtin                       varchar(128) not null default '',
    search_keywords            varchar(128) not null default '',
    availability               varchar(128) not null default '',
    availability_description   varchar(2048) not null default '',
    gift_wrapping_options_type varchar(128) not null default '',
    sort_order                 int8 not null default 0,
    condition                  varchar(2048) not null default '',
    is_condition_shown         bool not null default false,
    order_quantity_minimum     int8 not null default 0,
    order_quantity_maximum     int8 not null default 0,
    page_title                 varchar(512) not null default '',
    meta_description           varchar(1024) not null default '',
    date_created               varchar(128) not null default '',
    date_modified              varchar(128) not null default '',
    view_count                 int8 not null default 0,
    preorder_release_date      varchar(128) not null default '',
    preorder_message           varchar(2048) not null default '',
    is_preorder_only           bool not null default false,
    is_price_hidden            bool not null default false,
    price_hidden_label         varchar(2048) not null default '',
    base_variant_id            int8 not null default 0,
    custom_url                 jsonb
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists products;
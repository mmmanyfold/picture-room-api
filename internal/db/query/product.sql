-- name: CreateProduct :one
insert into products (id,
                      created_at,
                      updated_at,
                      deleted_at,
                      name,
                      type,
                      sku,
                      description,
                      weight,
                      Depth,
                      Height,
                      Price,
                      CostPrice,
                      RetailPrice,
                      SalePrice,
                      MapPrice,
                      tax_class_id,
                      product_tax_code,
                      calculated_price,
                      categories,
                      brand_id,
                      option_set_id,
                      option_set_display,
                      inventory_level,
                      inventory_warning_level,
                      inventory_tracking,
                      reviews_rating_sum,
                      reviews_count,
                      total_sold,
                      fixed_cost_shipping_price,
                      is_free_shipping,
                      is_visible,
                      is_featured,
                      warranty,
                      bin_picking_number,
                      layout_file,
                      upc,
                      mpn,
                      gtin,
                      search_keywords,
                      availability,
                      availability_description,
                      gift_wrapping_options_type,
                      sort_order,
                      condition,
                      is_condition_shown,
                      order_quantity_minimum,
                      order_quantity_maximum,
                      page_title,
                      meta_description,
                      date_created,
                      date_modified,
                      view_count,
                      preorder_release_date,
                      preorder_message,
                      is_preorder_only,
                      is_price_hidden,
                      price_hidden_label,
                      base_variant_id,
                      custom_url,
                      images,
                      variants)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
        $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46,
        $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62)

on conflict (id) do update set updated_at                 = now(),
                               name                       = $5,
                               type                       = $6,
                               sku                        = $7,
                               description                = $8,
                               weight                     = $9,
                               Depth                      = $10,
                               Height                     = $11,
                               Price                      = $12,
                               CostPrice                  = $13,
                               RetailPrice                = $14,
                               SalePrice                  = $15,
                               MapPrice                   = $16,
                               tax_class_id               = $17,
                               product_tax_code           = $18,
                               calculated_price           = $19,
                               categories                 = $20,
                               brand_id                   = $21,
                               option_set_id              = $22,
                               option_set_display         = $23,
                               inventory_level            = $24,
                               inventory_warning_level    = $25,
                               inventory_tracking         = $26,
                               reviews_rating_sum         = $27,
                               reviews_count              = $28,
                               total_sold                 = $29,
                               fixed_cost_shipping_price  = $30,
                               is_free_shipping           = $31,
                               is_visible                 = $32,
                               is_featured                = $33,
                               warranty                   = $34,
                               bin_picking_number         = $35,
                               layout_file                = $36,
                               upc                        = $37,
                               mpn                        = $38,
                               gtin                       = $39,
                               search_keywords            = $40,
                               availability               = $41,
                               availability_description   = $42,
                               gift_wrapping_options_type = $43,
                               sort_order                 = $44,
                               condition                  = $45,
                               is_condition_shown         = $46,
                               order_quantity_minimum     = $47,
                               order_quantity_maximum     = $48,
                               page_title                 = $49,
                               meta_description           = $50,
                               date_created               = $51,
                               date_modified              = $52,
                               view_count                 = $53,
                               preorder_release_date      = $54,
                               preorder_message           = $55,
                               is_preorder_only           = $56,
                               is_price_hidden            = $57,
                               price_hidden_label         = $58,
                               base_variant_id            = $59,
                               custom_url                 = $60,
                               images                     = $61,
                               variants                   = $62
returning *;

-- name: GetProduct :one
select
    products.*,
    brands.name as brand_name
from
    products
    inner join brands on products.brand_id = brands.id
where
    products.id = $1;

-- name: ListProducts :many
select
    products.*,
    brands.name as brand_name
from
    products
    inner join brands on products.brand_id = brands.id
where
    (@cat1::bigint = any (categories) or
    @cat2::bigint = any (categories) or
    @cat3::bigint = any (categories)) and
    (products.price between @minPrice::bigint and @maxPrice::bigint) and
    products.is_visible
order by
    (case when products.is_featured then 1 end) asc, products.id
limit @lim::bigint offset @off::bigint;

-- name: DeleteProduct :exec
delete
from
    products
where
    id = $1;

-- name: UpdateInventoryLevel :exec
update
    products
set
    inventory_level = $2
where
    id = $1;

-- name: ProductCount :one
select
    count(*)::bigint as featured_total,
    (select count(*) from products inner join brands on products.brand_id = brands.id where products.is_visible)::bigint as products_total
from
    products inner join brands on products.brand_id = brands.id
where
    (@cat1::bigint = any (categories) or
    @cat2::bigint = any (categories) or
    @cat3::bigint = any (categories)) and
    (products.price between @minPrice::bigint and @maxPrice::bigint) and
    products.is_visible;

prepare column_names (text) as
    select column_name from information_schema.columns where table_name = $1;

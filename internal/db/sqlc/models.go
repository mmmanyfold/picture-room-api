// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Brand struct {
	ID              int64           `json:"id"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       sql.NullTime    `json:"deleted_at"`
	Name            string          `json:"name"`
	PageTitle       string          `json:"page_title"`
	MetaKeywords    []string        `json:"meta_keywords"`
	MetaDescription string          `json:"meta_description"`
	SearchKeywords  string          `json:"search_keywords"`
	CustomUrl       json.RawMessage `json:"custom_url"`
}

type Product struct {
	ID                      int64           `json:"id"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
	DeletedAt               sql.NullTime    `json:"deleted_at"`
	Name                    string          `json:"name"`
	Type                    string          `json:"type"`
	Sku                     string          `json:"sku"`
	Description             string          `json:"description"`
	Weight                  float64         `json:"weight"`
	Depth                   float64         `json:"depth"`
	Height                  float64         `json:"height"`
	Price                   float64         `json:"price"`
	Costprice               float64         `json:"costprice"`
	Retailprice             float64         `json:"retailprice"`
	Saleprice               float64         `json:"saleprice"`
	Mapprice                float64         `json:"mapprice"`
	TaxClassID              int64           `json:"tax_class_id"`
	ProductTaxCode          string          `json:"product_tax_code"`
	CalculatedPrice         float64         `json:"calculated_price"`
	Categories              []int64         `json:"categories"`
	BrandID                 int64           `json:"brand_id"`
	OptionSetID             int64           `json:"option_set_id"`
	OptionSetDisplay        string          `json:"option_set_display"`
	InventoryLevel          int64           `json:"inventory_level"`
	InventoryWarningLevel   int64           `json:"inventory_warning_level"`
	InventoryTracking       string          `json:"inventory_tracking"`
	ReviewsRatingSum        int64           `json:"reviews_rating_sum"`
	ReviewsCount            int64           `json:"reviews_count"`
	TotalSold               int64           `json:"total_sold"`
	FixedCostShippingPrice  float64         `json:"fixed_cost_shipping_price"`
	IsFreeShipping          bool            `json:"is_free_shipping"`
	IsVisible               bool            `json:"is_visible"`
	IsFeatured              bool            `json:"is_featured"`
	Warranty                string          `json:"warranty"`
	BinPickingNumber        string          `json:"bin_picking_number"`
	LayoutFile              string          `json:"layout_file"`
	Upc                     string          `json:"upc"`
	Mpn                     string          `json:"mpn"`
	Gtin                    string          `json:"gtin"`
	SearchKeywords          string          `json:"search_keywords"`
	Availability            string          `json:"availability"`
	AvailabilityDescription string          `json:"availability_description"`
	GiftWrappingOptionsType string          `json:"gift_wrapping_options_type"`
	SortOrder               int64           `json:"sort_order"`
	Condition               string          `json:"condition"`
	IsConditionShown        bool            `json:"is_condition_shown"`
	OrderQuantityMinimum    int64           `json:"order_quantity_minimum"`
	OrderQuantityMaximum    int64           `json:"order_quantity_maximum"`
	PageTitle               string          `json:"page_title"`
	MetaDescription         string          `json:"meta_description"`
	DateCreated             string          `json:"date_created"`
	DateModified            string          `json:"date_modified"`
	ViewCount               int64           `json:"view_count"`
	PreorderReleaseDate     string          `json:"preorder_release_date"`
	PreorderMessage         string          `json:"preorder_message"`
	IsPreorderOnly          bool            `json:"is_preorder_only"`
	IsPriceHidden           bool            `json:"is_price_hidden"`
	PriceHiddenLabel        string          `json:"price_hidden_label"`
	BaseVariantID           int64           `json:"base_variant_id"`
	CustomUrl               json.RawMessage `json:"custom_url"`
	Images                  json.RawMessage `json:"images"`
	Variants                json.RawMessage `json:"variants"`
}
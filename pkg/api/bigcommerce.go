package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const includeOptions = "include=images,variants"

type Product struct {
	ID                      int64           `json:"id"`
	Name                    string          `json:"name"`
	Type                    string          `json:"type"`
	Sku                     string          `json:"sku"`
	Description             string          `json:"description"`
	Weight                  float64         `json:"weight"`
	Width                   float64         `json:"width"`
	Depth                   float64         `json:"depth"`
	Height                  float64         `json:"height"`
	Price                   float64         `json:"price"`
	CostPrice               float64         `json:"cost_price"`
	RetailPrice             float64         `json:"retail_price"`
	SalePrice               float64         `json:"sale_price"`
	MapPrice                float64         `json:"map_price"`
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
	PreorderReleaseDate     string          `json:"prorder_release_date"`
	PreorderMessage         string          `json:"preorder_message"`
	IsPreorderOnly          bool            `json:"is_preorder_only"`
	IsPriceHidden           bool            `json:"is_price_hidden"`
	PriceHiddenLabel        string          `json:"price_hidden_label"`
	CustomURL               json.RawMessage `json:"custom_url"`
	BaseVariantID           int64           `json:"base_variant_id"`
	Images                  json.RawMessage `json:"images"`
	Variants                json.RawMessage `json:"variants"`
	BrandName               string          `json:"brand_name,omitempty"`
}

type ProductData struct {
	Data Product `json: "data"`
}

type Brand struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	PageTitle       string          `json:"page_title"`
	MetaKeywords    []string        `json:"meta_keywords"`
	MetaDescription string          `json:"meta_description"`
	ImageURL        string          `json:"image_url"`
	SearchKeywords  string          `json:"search_keywords"`
	CustomURL       json.RawMessage `json:"custom_url"`
}

type BrandsPayload struct {
	Data []Brand `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			Links       struct {
				Next    string `json:"next"`
				Current string `json:"current"`
			} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

type ProductsPayload struct {
	Data []Product `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			Links       struct {
				Previous string `json:"previous"`
				Current  string `json:"current"`
			} `json:"links"`
			TooMany bool `json:"too_many"`
		} `json:"pagination"`
	} `json:"meta"`
}

type ProductWebhookPayload struct {
	CreatedAt int    `json:"created_at"`
	StoreID   string `json:"store_id"`
	Producer  string `json:"producer"`
	Scope     string `json:"scope"`
	Hash      string `json:"hash"`
	Data      struct {
		ID        int64  `json:"id"`
		Type      string `json:"type"`
		Inventory struct {
			ProductID int    `json:"product_id"`
			Method    string `json:"method"`
			Value     int64  `json:"value"`
		} `json:"inventory,omitempty"`
	} `json:"data"`
}

// FetchBrandsFromBC calls bigcommerce API to retrieve a list of brands which we use as artist pages
// https://api.bigcommerce.com/stores/{store-hash}/v3/catalog/brands
func FetchBrandsFromBC(bcConfigMap map[string]string, page, limit int) (BrandsPayload, error) {
	var brands BrandsPayload

	bcEndpoint := fmt.Sprintf("https://api.bigcommerce.com/stores/%s/v3/catalog/brands?limit=%d&page=%d", bcConfigMap["store"], limit, page)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, bcEndpoint, nil)
	if err != nil {
		return brands, fmt.Errorf("failed to make request for brand list, err: %w", err)
	}

	req.Header.Add("X-Auth-Client", bcConfigMap["client"])
	req.Header.Add("X-Auth-Token", bcConfigMap["token"])

	resp, err := client.Do(req)
	if err != nil {
		return brands, fmt.Errorf("failed to retrieve brand list from bigcommerce API, err: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&brands); err != nil {
		return brands, fmt.Errorf("error decoding json payload from bigcommerce API, err: %w", err)
	}

	if brands.Meta.Pagination.CurrentPage != brands.Meta.Pagination.TotalPages {
		nextPage := page + 1
		nextBrands, err := FetchBrandsFromBC(bcConfigMap, nextPage, limit)
		if err != nil {
			return brands, fmt.Errorf("error failed to get paginated set of brands from bigcommerce API, err: %w", err)
		}
		// update metadata
		brands.Meta = nextBrands.Meta
		// append brand records
		brands.Data = append(brands.Data, nextBrands.Data...)
	}
	return brands, nil
}

// FetchBrandFromBC calls bigcommerce API to retrieve a single of brand
// https://api.bigcommerce.com/stores/{store-hash}/v3/catalog/brands/{id}
func FetchBrandFromBC(bcConfigMap map[string]string, brandId int) (BrandsPayload, error) {
	var brand BrandsPayload

	bcEndpoint := fmt.Sprintf("https://api.bigcommerce.com/stores/%s/v3/catalog/brands/%d", bcConfigMap["store"], brandId)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, bcEndpoint, nil)
	if err != nil {
		return brand, fmt.Errorf("failed to make request for brand list, err: %w", err)
	}

	req.Header.Add("X-Auth-Client", bcConfigMap["client"])
	req.Header.Add("X-Auth-Token", bcConfigMap["token"])

	resp, err := client.Do(req)
	if err != nil {
		return brand, fmt.Errorf("failed to retrieve brand from bigcommerce API, err: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&brand); err != nil {
		return brand, fmt.Errorf("error decoding json payload from bigcommerce API, err: %w", err)
	}

	return brand, nil
}

// FetchProductsFromBC calls bigcommerce API to retrieve a list of products
// https://api.bigcommerce.com/stores/{store_hash}/v3/catalog/products
func FetchProductsFromBC(bcConfigMap map[string]string, page, limit int) (ProductsPayload, error) {
	var products ProductsPayload

	bcEndpoint := fmt.Sprintf("https://api.bigcommerce.com/stores/%s/v3/catalog/products?limit=%d&page=%d&%s",
		bcConfigMap["store"],
		limit,
		page,
		includeOptions,
	)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, bcEndpoint, nil)
	if err != nil {
		return ProductsPayload{}, fmt.Errorf("failed to make request for product list, err: %w", err)
	}

	req.Header.Add("X-Auth-Client", bcConfigMap["client"])
	req.Header.Add("X-Auth-Token", bcConfigMap["token"])

	resp, err := client.Do(req)
	if err != nil {
		return ProductsPayload{}, fmt.Errorf("failed to retrieve product list from bigcommerce API, err: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return ProductsPayload{}, fmt.Errorf("error decoding json payload from bigcommerce API, err: %w", err)
	}

	if products.Meta.Pagination.CurrentPage != products.Meta.Pagination.TotalPages {
		nextPage := page + 1
		nextProducts, err := FetchProductsFromBC(bcConfigMap, nextPage, limit)
		if err != nil {
			return ProductsPayload{}, fmt.Errorf("error failed to get paginated set of products from bigcommerce API, err: %w", err)
		}
		// update metadata
		products.Meta = nextProducts.Meta
		// append product records
		products.Data = append(products.Data, nextProducts.Data...)
	}

	return products, nil
}

// FetchProductFromBC calls bigcommerce API to retrieve a single product
// https://api.bigcommerce.com/stores/{store_hash}/v3/catalog/products/{productId}
func FetchProductFromBC(bcConfigMap map[string]string, productId int64) (Product, error) {
	var productData ProductData

	bcEndpoint := fmt.Sprintf("https://api.bigcommerce.com/stores/%s/v3/catalog/products/%d?%s",
		bcConfigMap["store"],
		productId,
		includeOptions,
	)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, bcEndpoint, nil)
	if err != nil {
		return Product{}, fmt.Errorf("failed to make request for single product, err: %w", err)
	}

	req.Header.Add("X-Auth-Client", bcConfigMap["client"])
	req.Header.Add("X-Auth-Token", bcConfigMap["token"])

	resp, err := client.Do(req)
	if err != nil {
		return Product{}, fmt.Errorf("failed to retrieve single product from bigcommerce API, productId: %d, err: %w", productId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return Product{}, MissingProductError
	}

	if err := json.NewDecoder(resp.Body).Decode(&productData); err != nil {
		return Product{}, fmt.Errorf("error decoding json payload from bigcommerce API, err: %w", err)
	}

	product := Product(productData.Data)
	return product, nil
}

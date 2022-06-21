package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	db "github.com/mmmanyfold/picture-room-api/internal/db/sqlc"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
)

var MissingProductError = errors.New("missing product from bigcommerce API")

type API struct {
	*db.Queries
	config *viper.Viper
	db     *sql.DB
}

type ProductResponse struct {
	Data db.GetProductRow `json:"data"`
}

type ProductsResponse struct {
	Data []db.ListProductsRow `json:"data"`
	Meta Meta                 `json:"meta"`
}

type Pagination struct {
	ProductsTotal int64 `json:"productsTotal"`
	FilteredTotal int64 `json:"filteredTotal"`
	ResultsTotal  int   `json:"resultsTotal"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

func NewAPI(c *viper.Viper, d *sql.DB) *API {
	return &API{
		db:      d,
		config:  c,
		Queries: db.New(d),
	}
}

func (a *API) GetBrands(w http.ResponseWriter, r *http.Request) {
	bcConfigMap := a.config.GetStringMapString("bigCommerce")
	if len(bcConfigMap) == 0 {
		http.Error(w, fmt.Sprintf("failed to retrieve brands list from bigcommerce API"), http.StatusInternalServerError)
		return
	}

	brandsData, err := FetchBrandsFromBC(bcConfigMap, 1, 100)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for brands list, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(brandsData)
}

func (a *API) GetBrand(w http.ResponseWriter, r *http.Request) {
	bcConfigMap := a.config.GetStringMapString("bigCommerce")
	if len(bcConfigMap) == 0 {
		http.Error(w, fmt.Sprintf("failed to retrieve brands list from bigcommerce API"), http.StatusInternalServerError)
		return
	}

	brandId, _ := strconv.Atoi(chi.URLParam(r, "id"))
	brands, err := FetchBrandFromBC(bcConfigMap, brandId)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for brands list, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(brands)
}

func (a *API) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	productId := chi.URLParam(r, "id")
	if productId == "" {
		http.Error(w, fmt.Sprintf("failed, product not found"), http.StatusNotFound)
		return
	}

	id, _ := strconv.Atoi(productId)
	product, err := a.Queries.GetProduct(ctx, int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for product, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	productsPayload := &ProductResponse{
		Data: product,
	}

	json.NewEncoder(w).Encode(productsPayload)
}

func (a *API) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	categories := r.URL.Query().Get("categories")
	minPrice, _ := strconv.ParseInt(r.URL.Query().Get("minPrice"), 10, 64)
	maxPrice, _ := strconv.ParseInt(r.URL.Query().Get("maxPrice"), 10, 64)

	cat, err := parseCategories(categories)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for product list, err: %v", err), http.StatusInternalServerError)
		return
	}

	params := db.ListProductsParams{
		Lim:      limit,
		Off:      offset,
		Cat1:     cat[0],
		Cat2:     cat[1],
		Cat3:     cat[2],
		Minprice: minPrice,
		Maxprice: maxPrice,
	}

	products, err := a.Queries.ListProducts(ctx, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for product list, err: %v", err), http.StatusInternalServerError)
		return
	}

	total, err := a.Queries.ProductCount(ctx, db.ProductCountParams{
		Cat1:     cat[0],
		Cat2:     cat[1],
		Cat3:     cat[2],
		Minprice: minPrice,
		Maxprice: maxPrice,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to make request for product list, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	productsPayload := &ProductsResponse{
		Data: products,
		Meta: Meta{
			Pagination{
				ProductsTotal: total.ProductsTotal,
				FilteredTotal: total.FeaturedTotal,
				ResultsTotal:  len(products),
			}},
	}

	json.NewEncoder(w).Encode(productsPayload)
}

func (a *API) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var pwp ProductWebhookPayload

	bcConfigMap := a.config.GetStringMapString("bigCommerce")
	if len(bcConfigMap) == 0 {
		// TODO: add bugsnag alert
		http.Error(w, fmt.Sprintf("failed to retrieve product list from bigcommerce API"), http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&pwp); err != nil {
		// TODO: add bugsnag alert
		log.Printf("failed to decode json payload from bigcommerce API webhook, err: %v", err)
		http.Error(w, "error decoding json payload from bigcommerce API webhook", http.StatusInternalServerError)
		return
	}

	err := a.WebhookRoute(ctx, bcConfigMap, &pwp)
	if err != nil {
		if errors.Is(err, MissingProductError) {
			w.WriteHeader(http.StatusOK)
			return
		}

		// TODO: add bugsnag alert
		log.Printf("err: %v", err)
		http.Error(w, fmt.Sprintf("error updating product data from bigcommerce API webhook, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) WebhookRoute(ctx context.Context, bcConfigMap map[string]string, pwp *ProductWebhookPayload) error {
	productId := pwp.Data.ID
	inventoryValue := pwp.Data.Inventory.Value

	switch pwp.Scope {
	case "store/product/created":
		fallthrough
	case "store/product/updated":
		product, err := FetchProductFromBC(bcConfigMap, productId)
		if err != nil {
			return err
		}

		params := CastProductParams(product)

		_, err = a.Queries.CreateProduct(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to upsert product from BC API Webhook with productId: %d, err: %v", productId, err)
		}

		return nil
	case "store/product/deleted":
		if err := a.Queries.DeleteProduct(ctx, productId); err != nil {
			return fmt.Errorf("failed to delete product from BC API Webhook with productId: %d, err: %v", productId, err)
		}
	case "store/product/inventory/updated":
		var params db.UpdateInventoryLevelParams

		params.ID = productId
		params.InventoryLevel = inventoryValue

		if err := a.Queries.UpdateInventoryLevel(ctx, params); err != nil {
			return err
		}
	}

	return nil
}

func CastProductParams(p Product) db.CreateProductParams {
	return db.CreateProductParams{
		ID:                      p.ID,
		Name:                    p.Name,
		Type:                    p.Type,
		Sku:                     p.Sku,
		Description:             p.Description,
		Weight:                  p.Weight,
		Depth:                   p.Depth,
		Height:                  p.Height,
		Price:                   p.Price,
		Costprice:               p.CostPrice,
		Retailprice:             p.RetailPrice,
		Saleprice:               p.SalePrice,
		Mapprice:                p.MapPrice,
		TaxClassID:              p.TaxClassID,
		ProductTaxCode:          p.ProductTaxCode,
		CalculatedPrice:         p.CalculatedPrice,
		Categories:              p.Categories,
		BrandID:                 p.BrandID,
		OptionSetID:             p.OptionSetID,
		OptionSetDisplay:        p.OptionSetDisplay,
		InventoryLevel:          p.InventoryLevel,
		InventoryWarningLevel:   p.InventoryWarningLevel,
		InventoryTracking:       p.InventoryTracking,
		ReviewsRatingSum:        p.ReviewsRatingSum,
		ReviewsCount:            p.ReviewsCount,
		TotalSold:               p.TotalSold,
		FixedCostShippingPrice:  p.FixedCostShippingPrice,
		IsFreeShipping:          p.IsFreeShipping,
		IsVisible:               p.IsVisible,
		IsFeatured:              p.IsFeatured,
		Warranty:                p.Warranty,
		BinPickingNumber:        p.BinPickingNumber,
		LayoutFile:              p.LayoutFile,
		Upc:                     p.Upc,
		Mpn:                     p.Mpn,
		Gtin:                    p.Gtin,
		SearchKeywords:          p.SearchKeywords,
		Availability:            p.Availability,
		AvailabilityDescription: p.AvailabilityDescription,
		GiftWrappingOptionsType: p.GiftWrappingOptionsType,
		SortOrder:               p.SortOrder,
		Condition:               p.Condition,
		IsConditionShown:        p.IsConditionShown,
		OrderQuantityMinimum:    p.OrderQuantityMinimum,
		OrderQuantityMaximum:    p.OrderQuantityMaximum,
		PageTitle:               p.PageTitle,
		MetaDescription:         p.MetaDescription,
		DateCreated:             p.DateCreated,
		DateModified:            p.DateModified,
		ViewCount:               p.ViewCount,
		PreorderReleaseDate:     p.PreorderReleaseDate,
		PreorderMessage:         p.PreorderMessage,
		IsPreorderOnly:          p.IsPreorderOnly,
		IsPriceHidden:           p.IsPriceHidden,
		PriceHiddenLabel:        p.PriceHiddenLabel,
		BaseVariantID:           p.BaseVariantID,
		CustomUrl:               p.CustomURL,
		Images:                  p.Images,
		Variants:                p.Variants,
	}
}

func CastBrandParams(b Brand) db.CreateBrandParams {
	return db.CreateBrandParams{
		ID:              int64(b.ID),
		Name:            b.Name,
		SearchKeywords:  b.SearchKeywords,
		PageTitle:       b.PageTitle,
		MetaDescription: b.MetaDescription,
		CustomUrl:       b.CustomURL,
	}
}

func parseCategories(query string) ([]int64, error) {
	var cat1 int64
	var cat2 int64
	var cat3 int64

	switch query {
	case "24", "25", "26":
		singleCat, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			return nil, err
		}
		cat1 = singleCat
	case "24,26":
		cat1 = 24
		cat2 = 26
	case "24,25":
		cat1 = 24
		cat2 = 25
	case "25,26":
		cat1 = 25
		cat2 = 26
	default:
		cat1 = 24
		cat2 = 25
		cat3 = 26
	}

	return []int64{cat1, cat2, cat3}, nil
}

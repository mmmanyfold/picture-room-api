package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	"github.com/mmmanyfold/picture-room-api/internal/config"
	"github.com/mmmanyfold/picture-room-api/pkg/api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var port string

	port = os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	env := os.Getenv("ENV")
	configPath := fmt.Sprintf("config/%s.yaml", env)

	if _, err := os.Stat(configPath); err != nil {
		fmt.Errorf("no config.yaml found, err: %w", err)
	}

	c := config.Read(configPath, nil)

	version := c.GetString("version")
	connStr := c.GetString("connStr")

	log.Printf("version::%s", version)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}
	defer database.Close()

	go scheduleInventorySync(c)

	API := api.NewAPI(c, database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// basic CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// middleware setup
	r.Use(
		corsHandler.Handler,
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,                             // log api request calls
		middleware.StripSlashes,                       // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,                          // recover from panics without crashing server
		middleware.Timeout(3000*time.Millisecond),     // Stop processing after 3 seconds
	)

	// obligatory health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/artists", API.GetBrands)
		r.Get("/artists/{id}", API.GetBrand)
		r.Get("/products", API.GetProducts)
		r.Get("/products/{id}", API.GetProduct)
		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/product", API.UpdateProducts)
		})
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func scheduleInventorySync(c *viper.Viper) {
	ctx := context.Background()
	bcConfigMap := c.GetStringMapString("bigCommerce")
	connStr := c.GetString("connStr")
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}
	defer database.Close()

	API := api.NewAPI(c, database)
	everyHour := time.NewTicker(1 * time.Hour)

	for {
		select {
		case <-everyHour.C:
			log.Println("worker::sync::products")
			products, err := api.FetchProductsFromBC(bcConfigMap, 1, 400)
			if err != nil {
				log.Println("worker::sync::products::fetch::failed")
				log.Fatalf("failed to retrieve products from BC api, err: %v\n", err)
			}

			if len(products.Data) > 0 {

				for i := 0; i < len(products.Data); i++ {
					params := api.CastProductParams(products.Data[i])
					_, err := API.Queries.CreateProduct(ctx, params)
					if err != nil {
						log.Println("worker::sync::products::save::failed")
						log.Fatalf("failed to persist products from BC api, err: %v\n", err)
					}
				}

				log.Println("worker::sync::products::save::successful")
			}

			log.Println("worker::sync::brands")

			brands, err := api.FetchBrandsFromBC(bcConfigMap, 1, 100)
			if err != nil {
				log.Println("worker::sync::brands::fetch::failed")
				log.Fatalf("failed to retrieve brands from BC api, err: %v\n", err)
			}

			if len(brands.Data) > 0 {

				for i := 0; i < len(brands.Data); i++ {
					params := api.CastBrandParams(brands.Data[i])
					_, err := API.Queries.CreateBrand(ctx, params)
					if err != nil {
						log.Println("worker::sync::brands::save::failed")
						log.Fatalf("failed to persist brands from BC api, err: %v\n", err)
					}
				}

				log.Println("worker::sync::brands::save::successful")
			}
		}
	}
}

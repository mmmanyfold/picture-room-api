// +build integration

package api

import (
	"github.com/go-chi/chi"
	"testing"
)

func NewTestRouter(api *API) *chi.Mux {
	r := chi.NewRouter()
	//r.Post("/webhook/slack", api.SlackHandler)
	return r
}

func Test_API(t *testing.T) {
	config := config.Read("../../config.yaml", nil)
	connStr := "postgres://postgres:postgres@127.0.0.1:5432/picture-room_test?sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}

	testAPI := NewAPI(config, db)
	router := NewTestRouter(testAPI)

	t.Cleanup(func() {
		ts := httptest.NewServer(router)
		defer ts.Close()
	})
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	orderservicev1 "github.com/sgrumley/kart-challenge/internal/services/order/v1"
	productservicev1 "github.com/sgrumley/kart-challenge/internal/services/product/v1"
	"github.com/sgrumley/kart-challenge/internal/store"
	"github.com/sgrumley/kart-challenge/pkg/idempotency"
	"github.com/sgrumley/kart-challenge/pkg/middleware"
	"github.com/sgrumley/kart-challenge/pkg/web"
)

func NewHandler(ctx context.Context, log slog.Logger, store *sqlx.DB) http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(time.Second * 10))

	router.Use(middleware.AddLogger(log))

	registerRoutes(router, store)

	return router
}

func registerRoutes(router *chi.Mux, client *sqlx.DB) {
	dbstore := store.New(client)
	idempotencyStore := idempotency.NewStore()

	routerv1 := chi.NewRouter()

	/*************************** PRODUCT ENDPOINTS ***************************/
	productService := productservicev1.NewService(dbstore)
	productService.GetRoutes(routerv1)

	/*************************** ORDER ENDPOINTS ***************************/
	orderService := orderservicev1.NewService(dbstore, idempotencyStore)
	orderService.GetRoutes(routerv1)

	router.Mount("/api/v1", routerv1)

	/*************************** HEALTHCHECK  ***************************/
	router.Post("/health", healthCheck)

	// list all the available endpoints
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%-6s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	web.RespondNoContent(w)
}

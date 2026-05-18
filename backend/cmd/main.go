package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cubancodepath/zerobudget/backend/internal/accounts"
	"github.com/cubancodepath/zerobudget/backend/internal/categories"
	"github.com/cubancodepath/zerobudget/backend/internal/payees"
	"github.com/cubancodepath/zerobudget/backend/internal/transactions"
	db "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	queries := db.New(pool)
	accountsService := accounts.NewService(queries)
	accountsHandler := accounts.NewHandler(accountsService)
	categoriesService := categories.NewService(queries)
	categoriesHandler := categories.NewHandler(categoriesService)
	payeesService := payees.NewService(queries)
	payeesHandler := payees.NewHandler(payeesService)
	transactionsService := transactions.NewService(pool, queries)
	transactionsHandler := transactions.NewHandler(transactionsService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	accountsHandler.RegisterRoutes(mux)
	categoriesHandler.RegisterRoutes(mux)
	payeesHandler.RegisterRoutes(mux)
	transactionsHandler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux, allowedOrigin)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func withCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Vary", "Access-Control-Request-Method")
		w.Header().Set("Vary", "Access-Control-Request-Headers")

		origin := r.Header.Get("Origin")
		if origin != "" && origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

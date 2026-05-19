package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cubancodepath/zerobudget/backend/internal/accounts"
	db "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/categories"
	"github.com/cubancodepath/zerobudget/backend/internal/env"
	"github.com/cubancodepath/zerobudget/backend/internal/payees"
	"github.com/cubancodepath/zerobudget/backend/internal/transactions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config
	db     *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn: env.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/zerobudget?sslmode=disable"),
		},
	}

	pool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	app := &application{
		config: cfg,
		db:     pool,
	}

	queries := db.New(app.db)
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

	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: withCORS(mux, allowedOrigin),
	}

	go func() {
		log.Printf("server listening on %s", app.config.addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced server shutdown: %v", err)
	}

	log.Println("server shutdown gracefully")
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

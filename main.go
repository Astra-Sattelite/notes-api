package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Getting .env creds
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbCreds := os.Getenv("DB_CREDS")
	port := os.Getenv("PORT")

	// Initialize PostgreSQL connection pool
	db, err := pgxpool.New(ctx, dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify database connection is active
	if err := db.Ping(ctx); err != nil {
		log.Fatal("postgres: ", err)
	}
	log.Println("postgres connected")

	// Connecting to redis
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis: ", err)
	}
	log.Println("redis connected")

	router := chi.NewRouter()

	router.Get("/ping", func(resp http.ResponseWriter, req *http.Request) {
		var now time.Time
		// Check if data available
		err := db.QueryRow(req.Context(), "SELECT now()").Scan(&now)
		if err != nil {
			http.Error(resp, "db error", http.StatusInternalServerError)
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(map[string]any{"status": "ok", "db_time": now})
	})

	router.Get("/redis-test", func(resp http.ResponseWriter, req *http.Request) {
		// Write test key in redis
		err := rdb.Set(req.Context(), "hello", "world", 30*time.Second).Err()
		if err != nil {
			http.Error(resp, "redis error", http.StatusInternalServerError)
			return
		}

		// Trying to read test key
		val, err := rdb.Get(req.Context(), "hello").Result()
		if err != nil {
			http.Error(resp, "redis error", http.StatusInternalServerError)
			return
		}

		// Return test key
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(map[string]string{"hello": val})
	})

	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}

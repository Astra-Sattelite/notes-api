package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"notes-api/rcache"
	"notes-api/health"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbCreds := os.Getenv("DB_CREDS")
	port := os.Getenv("PORT")

	db, err := pgxpool.New(ctx, dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal("postgres: ", err)
	}
	log.Println("postgres connected")

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis: ", err)
	}
	log.Println("redis connected")

	router := chi.NewRouter()

	health.RegisterRoutes(router, db)
	rcache.RegisterRoutes(router, rdb)

	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
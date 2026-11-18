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
	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}

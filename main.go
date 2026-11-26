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
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"notes-api/docs"
	"notes-api/health"
	"notes-api/rcache"
)

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateNoteRequest struct {
	Title string `json:"title" example:"First note"`
	Body  string `json:"body"  example:"Hello from Swagger"`
}

type Server struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

// @title   Notes API
// @version 1.0
// @BasePath /
func main() {
	ctx := context.Background()

	// .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbCreds := os.Getenv("DB_CREDS")
	port := os.Getenv("PORT")
	rdbCreds := os.Getenv("RDB_CREDS")
	if dbCreds == "" || port == "" || rdbCreds == "" {
		log.Fatal("DB_CREDS, RDB_CREDS and PORT must be set in .env")
	}

	// Swagger
	if value := os.Getenv("SWAGGER_HOST"); value != "" {
		docs.SwaggerInfo.Host = value
	}

	// Postgres
	db, err := pgxpool.New(ctx, dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal("postgres: ", err)
	}
	log.Println("postgres connected")

	// Redis
	rdb := redis.NewClient(&redis.Options{Addr: rdbCreds})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis: ", err)
	}
	log.Println("redis connected")

	s := &Server{db: db, rdb: rdb}

	router := chi.NewRouter()
	router.Post("/notes", s.createNote)

	health.RegisterRoutes(router, db)
	rcache.RegisterRoutes(router, rdb)

	if os.Getenv("APP_ENV") == "dev" {
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
		log.Println("swagger: http://localhost:" + port + "/swagger/index.html")
	}

	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// createNote godoc
// @Summary     Create a note
// @Description Creates a note and returns it with id and created_at
// @Tags        notes
// @Accept      json
// @Produce     json
// @Param       note body CreateNoteRequest true "Note data"
// @Success     201 {object} Note
// @Failure     400 {string} string "bad json or title is required"
// @Failure     500 {string} string "internal error"
// @Router      /notes [post]
func (s *Server) createNote(resp http.ResponseWriter, req *http.Request) {
	// Read JSON from body
	var in CreateNoteRequest

	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(resp, "bad json", http.StatusBadRequest)
		return
	}
	if in.Title == "" {
		http.Error(resp, "title is required", http.StatusBadRequest)
		return
	}

	// INSERT note-data into Postgres, and getting it back
	var note Note
	err := s.db.QueryRow(req.Context(),
		`
			INSERT INTO notes (title, body)
			VALUES ($1, $2)
			RETURNING id, title, body, created_at
		`,
		in.Title, in.Body,
	).Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		log.Println("insert:", err)
		http.Error(resp, "internal error", http.StatusInternalServerError)
		return
	}

	// Giving back Note as JSON
	writeJSON(resp, http.StatusCreated, note)
}

func writeJSON(resp http.ResponseWriter, status int, value any) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	json.NewEncoder(resp).Encode(value)
}

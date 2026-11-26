package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func RegisterRoutes(r chi.Router, db *pgxpool.Pool) {
	h := &Handler{db: db}
	r.Get("/ping", h.ping)
}

func (h *Handler) ping(resp http.ResponseWriter, req *http.Request) {
	var now time.Time
	err := h.db.QueryRow(req.Context(), "SELECT now()").Scan(&now)
	if err != nil {
		http.Error(resp, "db error", http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(map[string]any{"status": "ok", "db_time": now})
}

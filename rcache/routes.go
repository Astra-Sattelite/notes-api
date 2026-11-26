package rcache

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	rdb *redis.Client
}

func RegisterRoutes(r chi.Router, rdb *redis.Client) {
	h := &Handler{rdb: rdb}
	r.Get("/redis-test", h.test)
}

func (h *Handler) test(resp http.ResponseWriter, req *http.Request) {
	err := h.rdb.Set(req.Context(), "hello", "world", 30*time.Second).Err()
	if err != nil {
		http.Error(resp, "redis error", http.StatusInternalServerError)
		return
	}

	val, err := h.rdb.Get(req.Context(), "hello").Result()
	if err != nil {
		http.Error(resp, "redis error", http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(map[string]string{"hello": val})
}

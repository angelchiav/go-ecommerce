package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/angelchiav/go-ecommerce/internal/httpx"
)

type Health struct {
	db *sql.DB
}

func NewHealth(db *sql.DB) *Health { return &Health{db: db} }

func (h *Health) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		httpx.Error(w, http.StatusServiceUnavailable, "db_down")
		return
	}

	var one int
	if err := h.db.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
		httpx.Error(w, http.StatusServiceUnavailable, "db_query_failed")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

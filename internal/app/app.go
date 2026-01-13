package app

import (
	"net/http"

	"github.com/angelchiav/go-ecommerce/internal/config"
	"github.com/angelchiav/go-ecommerce/internal/db"
	"github.com/angelchiav/go-ecommerce/internal/handlers"
	"github.com/angelchiav/go-ecommerce/internal/httpx"
)

type App struct {
	handler http.Handler
}

func New(cfg config.Config) (*App, error) {
	conn, err := db.OpenPostgres(cfg.DBURL)
	if err != nil {
		return nil, err
	}

	r := httpx.NewRouter()
	health := handlers.NewHealth(conn)

	r.Handle("GET", "/health", health.Get)

	h := httpx.Recover(httpx.Logger(r))

	return &App{handler: h}, nil
}

func (a *App) Handler() http.Handler { return a.handler }

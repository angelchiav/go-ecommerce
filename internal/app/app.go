package app

import (
	"net/http"

	"github.com/angelchiav/go-ecommerce/internal/config"
	"github.com/angelchiav/go-ecommerce/internal/db"
	"github.com/angelchiav/go-ecommerce/internal/handlers"
	"github.com/angelchiav/go-ecommerce/internal/httpx"
	"github.com/angelchiav/go-ecommerce/internal/service"
	"github.com/angelchiav/go-ecommerce/internal/sqlc"
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

	q := sqlc.New(conn)

	cartSvc := service.NewCartService(conn, q)
	cartH := handlers.NewCart(cartSvc)

	r.Handle("GET", "/v1/cart", cartH.Get)
	r.Handle("GET", "/health", health.Get)
	r.Handle("POST", "/v1/cart/items", cartH.AddItem)
	r.Handle("PATCH", "/v1/cart/items/{id}", cartH.UpdateItemQty)
	r.Handle("DELETE", "/v1/cart/items/{id}", cartH.DeleteItem)

	h := httpx.Recover(httpx.Logger(r))

	return &App{handler: h}, nil
}

func (a *App) Handler() http.Handler { return a.handler }

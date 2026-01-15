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

	authSvc := service.NewAuthService(q, cfg.JWTSecret)
	authH := handlers.NewAuth(authSvc, q)
	authMW := httpx.AuthJWT(cfg.JWTSecret)

	// PUBLIC
	r.Handle("GET", "/health", health.Get)
	r.Handle("POST", "/v1/auth/register", authH.Register)
	r.Handle("POST", "/v1/auth/login", authH.Login)

	// PRIVATE
	r.Handle("GET", "/v1/me", authMW(authH.Me))
	r.Handle("GET", "/v1/cart", authMW(cartH.Get))
	r.Handle("POST", "/v1/cart/items", authMW(cartH.AddItem))
	r.Handle("PATCH", "/v1/cart/items/{id}", authMW(cartH.UpdateItemQty))
	r.Handle("DELETE", "/v1/cart/items/{id}", authMW(cartH.DeleteItem))

	h := httpx.Recover(httpx.Logger(r))

	return &App{handler: h}, nil
}

func (a *App) Handler() http.Handler { return a.handler }

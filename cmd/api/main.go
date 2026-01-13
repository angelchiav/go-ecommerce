package main

import (
	"log"
	"net/http"
	"time"

	"github.com/angelchiav/go-ecommerce/internal/app"
	"github.com/angelchiav/go-ecommerce/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	log.Printf("DB_URL:%q", cfg.DBURL)
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      a.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on: %s", cfg.Addr)
	log.Fatal(srv.ListenAndServe())
}

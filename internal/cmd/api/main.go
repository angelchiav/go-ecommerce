package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-ecommerce/internal/config"
)

func main() {
	cfg := config.Load()
	a, err := app.New()
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

package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr      string
	DBURL     string
	JWTSecret string
}

func Load() Config {
	_ = loadDotEnvFromModuleRoot()

	return Config{
		Addr:      env("ADDR", ":8080"),
		DBURL:     mustEnv("DB_URL"),
		JWTSecret: mustEnv("JWT_SECRET"),
	}
}

func loadDotEnvFromModuleRoot() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {

			_ = godotenv.Load(filepath.Join(dir, ".env"))
			return nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}

func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing required env var" + k)
	}
	return v
}

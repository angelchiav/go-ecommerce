package config

import "os"

type Config struct {
	Addr  string
	DBURL string
}

func Load() Config {
	return Config{
		Addr:  os.Getenv("Addr"),
		DBURL: os.Getenv("DBURL"),
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

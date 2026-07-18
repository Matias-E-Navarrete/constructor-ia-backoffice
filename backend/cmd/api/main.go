// Command api is the composition root's entry point: load config, build the
// wired router (see internal/platform/httpserver.Build), start listening.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"rimu/backend/internal/platform/config"
	"rimu/backend/internal/platform/httpserver"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	router, pool, err := httpserver.Build(context.Background(), cfg)
	if err != nil {
		log.Fatalf("build server: %v", err)
	}
	defer pool.Close()

	log.Printf("rimu backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}

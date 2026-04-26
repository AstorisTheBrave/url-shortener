package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port    := envOr("PORT",     "8080")
	baseURL := envOr("BASE_URL", "http://localhost:"+port)
	dbPath  := envOr("DB_PATH",  "urls.json")

	h := &handler{
		store:   NewStore(dbPath),
		baseURL: baseURL,
	}

	fmt.Printf("url-shortener on :%s  (base: %s)\n", port, baseURL)
	log.Fatal(http.ListenAndServe(":"+port, h))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

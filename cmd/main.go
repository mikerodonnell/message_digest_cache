package main

import (
	"log"
	"net/http"

	"github.com/mikerodonnell/message_digest_cache/pkg/api"
	"github.com/mikerodonnell/message_digest_cache/pkg/persist"
)

// important to use :8000, not localhost:8000, for docker ports to forward
const host = ":8000"

func main() {
	log.Println("connecting to distributed cache")
	cache, err := persist.NewRedisCache()
	if err != nil {
		log.Fatal("failed to connect to distributed cache", err)
	}

	// this isn't guaranteed to execute if process is killed; in a real implementation we'd have a /stop endpoint or similar
	// (though CLIENT LIST shows killing the process doesn't leak connections)
	defer cache.Close()

	log.Println("initializing API router")
	router := api.NewRouter(cache)

	log.Println("starting server at: ", host)
	log.Fatal(http.ListenAndServe(host, router))
}

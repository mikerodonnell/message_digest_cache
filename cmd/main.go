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
	log.Println("initializing local cache")
	localCache := persist.NewLocalCache()

	log.Println("connecting to distributed cache")
	distributedCache, err := persist.NewRedisCache()
	if err != nil {
		log.Fatal("failed to connect to distributed cache", err)
	}
	// this isn't guaranteed to execute if process is killed; in a real implementation we'd have a /stop endpoint or similar
	// (though CLIENT LIST shows killing the process doesn't leak connections)
	defer distributedCache.Close()

	log.Println("initializing API router")
	// pass localCache first to be our primary, and distributedCache second to be the backup
	// in a multi app node deployment (or if an app node is restarted), we can safely just go to redis
	// because the keys are always deterministic SHA digests of the values, there's no updates,
	// and therefore no synchronization to worry about between local digests on separate nodes
	router := api.NewRouter(localCache, distributedCache)

	log.Println("starting server at ", host)
	log.Fatal(http.ListenAndServe(host, router))
}

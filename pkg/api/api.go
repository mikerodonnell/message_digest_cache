package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikerodonnell/message_digest_cache/pkg/persist"
)

var caches []persist.Cache

type putRequest struct {
	Message string `json:"message"`
}

type getResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type putResponse struct {
	Digest string `json:"digest,omitempty"`
	Error  string `json:"error,omitempty"`
}

// NewRouter creates a mux.Router that uses the given cache(s) to
func NewRouter(newCaches ...persist.Cache) *mux.Router {
	caches = newCaches

	router := mux.NewRouter()

	router.HandleFunc("/messages", put).Methods("POST")

	router.HandleFunc("/messages/{digest}", get).Methods("GET")

	// these are just required to respond with 400 instead of 405 when {digest} is missing from request
	router.HandleFunc("/messages", get).Methods("GET")
	router.HandleFunc("/messages/", get).Methods("GET")

	return router
}

func get(w http.ResponseWriter, r *http.Request) {
	response := getResponse{}

	digest := mux.Vars(r)["digest"]

	if len(digest) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = "`digest` request variable is required"
		json.NewEncoder(w).Encode(response)

		return
	}

	for _, cache := range caches {
		message := cache.Get(digest)
		if len(message) > 1 {
			// cache hit! populate response and return
			response.Message = message
			json.NewEncoder(w).Encode(response)

			return
		}
	}

	// cache miss; no digest for this message
	w.WriteHeader(http.StatusNotFound)
	response.Message = "message not found"
	json.NewEncoder(w).Encode(response)

	return
}

func put(w http.ResponseWriter, r *http.Request) {
	response := putResponse{}

	decoder := json.NewDecoder(r.Body)
	var body putRequest
	err := decoder.Decode(&body)
	if err != nil {
		response.Error = "malformed request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)

		return
	}

	message := body.Message
	digestBytes := sha256.Sum256([]byte(message))
	digest := fmt.Sprintf("%x", digestBytes) // %x for lowercase hex characters

	// put in each cache
	for _, cache := range caches {
		err = cache.Put(digest, message)
		if err != nil {
			// TODO: in a production implementation we'd likely want transactionality here, across all caches

			sanitizedMessage := fmt.Sprintf("server error storing digest for %s", message)
			log.Println(sanitizedMessage, err)
			response.Error = sanitizedMessage
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)

			return
		}
	}

	response.Digest = digest
	json.NewEncoder(w).Encode(response)

	return
}

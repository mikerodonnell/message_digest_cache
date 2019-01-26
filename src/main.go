package main

import (
	"fmt"
	"log"
	"net/http"

	"api"
)

func main() {
	fmt.Println("starting server")

	router := api.NewRouter()

	log.Fatal(http.ListenAndServe(":8000", router))
}

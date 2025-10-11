package main

import (
	"log"
	"net/http"
)

func main() {
	mux := Mux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", dir, port)
	log.Fatal(server.ListenAndServe())
}

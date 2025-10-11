package main

import (
	"log"
	"net/http"
)

func main() {
	_ = http.NewServeMux()
	server := http.Server{}
	server.Addr = ":8080"
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

package main

import (
	"fmt"
	"net/http"
)

func handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.hitsReset()
	w.WriteHeader(http.StatusOK)
}

func handlerMetrics(w http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Hits: %v\n", cfg.getHits())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

package main

import (
	"fmt"
	"net/http"
	"os"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	w.Write([]byte("\n"))
}

func (cfg *apiCfg) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	req.Header.Set("Content-Type", "text/html")
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiCfg) handlerReset(w http.ResponseWriter, req *http.Request) {
	pf := os.Getenv("PLATFORM")
	if pf != "dev" {
		http.Error(w, "Access forbidden", http.StatusForbidden)
	}

	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	cfg.db.ResetDatabase(req.Context())
}

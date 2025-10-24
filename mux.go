package main

import "net/http"

func Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(dir))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/users", handlerUsers)
	mux.HandleFunc("POST /api/chirps", handlerChirps)
	mux.HandleFunc("GET /api/chirps", handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirp)
	mux.HandleFunc("POST /api/login", handlerLogin)
	mux.HandleFunc("POST /api/refresh", handlerRefresh)
	mux.HandleFunc("POST /api/revoke", handlerRevoke)

	return mux
}

package main

import "net/http"

func Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(dir))))
	mux.HandleFunc("/healthz/", handlerReadiness)
	mux.HandleFunc("/metrics", handlerMetrics)
	mux.HandleFunc("/reset", http.HandlerFunc(handlerReset))

	return mux
}

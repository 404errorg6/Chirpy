package main

import (
	"net/http"
	"sync/atomic"
)

var (
	dir  = http.Dir(".")
	port = "8080"
	cfg  = apiCfg{}
)

type apiCfg struct {
	fileServerHits atomic.Int32
}

func (cfg *apiCfg) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			cfg.fileServerHits.Add(1)
			next.ServeHTTP(w, req)
		},
	)
}

func (cfg *apiCfg) getHits() int32 {
	return cfg.fileServerHits.Load()
}

func (cfg *apiCfg) hitsReset() {
	cfg.fileServerHits.Store(0)
}

package main

import (
	"net/http"
	"sync/atomic"
	"time"

	"chirpy/internal/database"

	"github.com/google/uuid"
)

var (
	dir  = http.Dir(".")
	port = "8080"
	cfg  = apiCfg{}
)

type apiCfg struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	secret         string
}

type Users struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password,omitempty"`
	Token          string    `json:"token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiCfg) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			cfg.fileServerHits.Add(1)
			next.ServeHTTP(w, req)
		},
	)
}

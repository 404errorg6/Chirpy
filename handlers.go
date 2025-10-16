package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"chirpy/internal/database"

	"github.com/google/uuid"
)

func handlerChirps(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		respondErrJSON(w, 500, err)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		respondErrJSON(w, 400, err)
		return
	}
	if len(input.Body) > 140 {
		respondErrJSON(w, 400, fmt.Errorf("chirp is too long"))
		return
	}

	params := database.CreateChirpParams{
		Body:   input.Body,
		UserID: input.UserID,
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), params)
	if err != nil {
		respondErrJSON(w, 500, err)
		return
	}

	respondJSON(w, 201, chirp)
}

func handlerUsers(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		respondErrJSON(w, 500, err)
		return
	}
	if err := json.Unmarshal(body, &input); err != nil {
		respondErrJSON(w, 400, err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), input.Email)
	if err != nil {
		respondErrJSON(w, 500, err)
		return
	}
	user.Email = input.Email
	respondJSON(w, 200, user)
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
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
		respondErrJSON(w, 403, fmt.Errorf("Forbidden"))
	}

	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	cfg.db.ResetDatabase(req.Context())
}

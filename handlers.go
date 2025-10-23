package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"chirpy/internal/auth"
	"chirpy/internal/database"

	"github.com/google/uuid"
)

func handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	idstr := req.PathValue("chirpID")
	id, err := uuid.Parse(idstr)
	if err != nil {
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), id)
	if err != nil {
		http.Error(w, "Could not get chirp", 404)
		return
	}

	respondJSON(w, 200, chirp)
}

func handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		http.Error(w, "Could not get data", http.StatusNoContent)
		return
	}

	respondJSON(w, 200, chirps)
}

func handlerChirps(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Error occured while reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, "Error unmarhsalling data", http.StatusBadRequest)
		return
	}
	if len(input.Body) > 140 {
		http.Error(w, "Chirp is too long", http.StatusBadRequest)
		return
	}
	input.Body = cleanBody(input.Body)

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Could not get token from http header\n", http.StatusBadRequest)
		return
	}

	gotID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil || gotID != input.UserID {
		http.Error(w, "Unauthorized\n", 401)
		return
	}

	params := database.CreateChirpParams{
		Body:   input.Body,
		UserID: input.UserID,
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), params)
	if err != nil {
		er := fmt.Sprintf("Error occured while creating chirp: %v", err)
		http.Error(w, er, http.StatusInternalServerError)
		return
	}

	respondJSON(w, 201, chirp)
}

func handlerUsers(w http.ResponseWriter, req *http.Request) {
	var user Users
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Error occured while reading request body", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &input); err != nil {
		er := fmt.Sprintf("Error unmarhsalling data: %v", err)
		http.Error(w, er, http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		er := fmt.Sprintf("Error occured while hashing password: %v\n", err)
		http.Error(w, er, 500)
		return
	}

	params := database.CreateUserParams{
		Email:          input.Email,
		HashedPassword: hash,
	}

	tmpUser, err := cfg.db.CreateUser(req.Context(), params)
	if err != nil {
		er := fmt.Sprintf("Error occured while creating user: %v", err)
		http.Error(w, er, http.StatusInternalServerError)
		return
	}

	user.ID = tmpUser.ID
	user.Email = tmpUser.Email
	user.CreatedAt = tmpUser.CreatedAt
	user.UpdatedAt = tmpUser.UpdatedAt

	respondJSON(w, 200, user)
	user.HashedPassword = tmpUser.HashedPassword
}

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

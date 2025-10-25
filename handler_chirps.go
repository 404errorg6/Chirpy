package main

import (
	"fmt"
	"net/http"

	"chirpy/internal/auth"
	"chirpy/internal/database"

	"github.com/google/uuid"
)

func handlerDelChirp(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid chirpID", 400)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Could not get token from http header", 400)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		http.Error(w, "Invalid token", 401)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		http.Error(w, "Invalid chirpID", 401)
		return
	}

	if userID != chirp.UserID {
		http.Error(w, "Forbidden", 403)
		return
	}

	err = cfg.db.DelChirp(req.Context(), chirpID)
	if err != nil {
		http.Error(w, "Could not delete chirp", 404)
		return
	}

	w.WriteHeader(204)
}

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

	err := unmarshaller(req.Body, &input)
	req.Body.Close()

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

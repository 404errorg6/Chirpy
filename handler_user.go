package main

import (
	"fmt"
	"net/http"
	"time"

	"chirpy/internal/auth"
	"chirpy/internal/database"

	"github.com/google/uuid"
)

func handlerUpgradeSubs(w http.ResponseWriter, req *http.Request) {
	type DataT struct {
		UserID string `json:"user_id"`
	}
	var input struct {
		Event string `json:"event"`
		Data  DataT  `json:"data"`
	}

	api, err := auth.GetAPIKey(req.Header)
	if err != nil || api != cfg.polkaKey {
		http.Error(w, "Invalid api key", 401)
		return
	}

	err = unmarshaller(req.Body, &input)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Error unmarhsalling data", 400)
		return
	}

	idStr := input.Data.UserID
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid userID", 400)
		return
	}

	if input.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.db.GetUserByID(req.Context(), userID)
	if err != nil {
		http.Error(w, "Invalid userID", 400)
		return
	}

	err = cfg.db.UpgradeUser(req.Context(), userID)
	if err != nil {
		http.Error(w, "Could not upgrade user", 500)
		return
	}

	w.WriteHeader(204)
}

func handlerUpdateUsers(w http.ResponseWriter, req *http.Request) {
	var user Users
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Could not get token from header", 400)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		http.Error(w, "Unauthorized\n", 401)
		return
	}

	err = unmarshaller(req.Body, &input)
	req.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarhsalling data: %v", err), 401)
		return
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Could not hash the password", http.StatusInternalServerError)
	}

	params := database.UpdateUserParams{
		ID:             id,
		Email:          input.Email,
		HashedPassword: hash,
	}

	err = cfg.db.UpdateUser(req.Context(), params)
	if err != nil {
		http.Error(w, "Could not update user data", http.StatusInternalServerError)
		return
	}

	user.Email = input.Email
	user.ID = id
	user.UpdatedAt = time.Now()
	user.HashedPassword = ""

	respondJSON(w, 200, user)
	user.HashedPassword = hash
}

func handlerUsers(w http.ResponseWriter, req *http.Request) {
	var user Users

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := unmarshaller(req.Body, &input)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Error unmarshalling request body", 400)
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

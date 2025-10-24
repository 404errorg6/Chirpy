package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"chirpy/internal/auth"
	"chirpy/internal/database"
)

func handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Token not found in header\n", 401)
		return
	}

	params := database.UpdateRefreshTokenParams{
		Token:     token,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = cfg.db.UpdateRefreshToken(req.Context(), params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not revoke token: %v", err), 500)
		return
	}

	w.WriteHeader(204)
}

func handlerRefresh(w http.ResponseWriter, req *http.Request) {
	var output struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Token not found in header\n", 401)
		return
	}
	_, err = cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil {
		http.Error(w, "Invalid token", 401)
		return
	}

	output.Token = token
	respondJSON(w, 200, output)
}

func handlerLogin(w http.ResponseWriter, req *http.Request) {
	var user Users

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := unmarshaller(req.Body, &input)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Could not read the body\n", 400)
		return
	}

	tmpUser, err := cfg.db.GetUserByEmail(req.Context(), input.Email)
	if err != nil {
		http.Error(w, "Incorrect password or email", 401)
		return
	}

	match, err := auth.CheckPasswordHash(input.Password, tmpUser.HashedPassword)
	if err != nil {
		http.Error(w, "Unexpcted error occured", 500)
	}
	if !match {
		http.Error(w, "Incorrect password or email", 401)
		return
	}

	accessToken, err := auth.MakeJWT(tmpUser.ID, cfg.secret, 1*time.Hour)
	if err != nil {
		http.Error(w, "Cound not make JWT\n", 500)
		return
	}

	rToken, _ := auth.MakeRefreshToken()
	// Storing refreshToken in db
	params := database.CreateRefreshTokenParams{
		Token:     rToken,
		UpdatedAt: time.Now(),
		UserID:    tmpUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
	}

	refreshToken, err := cfg.db.CreateRefreshToken(req.Context(), params)
	if err != nil {
		http.Error(w, "Could not create refresh token", 500)
		return
	}

	user.Email = tmpUser.Email
	user.CreatedAt = tmpUser.CreatedAt
	user.ID = tmpUser.ID
	user.UpdatedAt = tmpUser.UpdatedAt
	user.Token = accessToken
	user.RefreshToken = refreshToken.Token

	respondJSON(w, 200, user)
}

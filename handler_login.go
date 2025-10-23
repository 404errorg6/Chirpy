package main

import (
	"fmt"
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func handlerLogin(w http.ResponseWriter, req *http.Request) {
	var user struct {
		Users
	}
	var input struct {
		Email    string        `json:"email"`
		Password string        `json:"password"`
		Expiry   time.Duration `json:"expiry_in_seconds,omitempty"`
	}

	err := unmarshaller(w, req.Body, &input)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Could not read the body", 400)
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

	if input.Expiry == 0 || input.Expiry > time.Hour {
		input.Expiry = time.Hour
	}

	token, err := auth.MakeJWT(tmpUser.ID, cfg.secret, input.Expiry)
	if err != nil {
		http.Error(w, "Cound not make JWT\n", 500)
		return
	}

	user.Email = tmpUser.Email
	user.CreatedAt = tmpUser.CreatedAt
	user.ID = tmpUser.ID
	user.UpdatedAt = tmpUser.UpdatedAt
	user.Token = token
	fmt.Printf("Token: %v\n", token)

	respondJSON(w, 200, user)
}

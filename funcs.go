package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func respondErrJSON(w http.ResponseWriter, code int, err error) {
	var errored struct {
		Error string `json:"error"`
	}
	errored.Error = fmt.Sprint(err)

	respondJSON(w, code, errored)
}

func respondJSON(w http.ResponseWriter, code int, body any) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(body)
	if err != nil {
		respondErrJSON(w, 500, err)
	}
	w.Write(data)
	w.Write([]byte("\n"))
}

func cleanBody(body string) string {
	notAllowed := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		if word == notAllowed[0] || word == notAllowed[1] || word == notAllowed[2] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

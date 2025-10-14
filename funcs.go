package main

import "strings"

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

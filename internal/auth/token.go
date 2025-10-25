package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("authorization")
	apiKey, ok := strings.CutPrefix(auth, "ApiKey ")
	if !ok || apiKey == "" {
		return apiKey, fmt.Errorf("Api key not found")
	}
	return apiKey, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("authorization")
	token, ok := strings.CutPrefix(bearer, "Bearer ")
	if !ok || token == "" {
		return token, fmt.Errorf("No token found")
	}
	return token, nil
}

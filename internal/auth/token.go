package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("authorization")
	token, ok := strings.CutPrefix(bearer, "Bearer ")
	if !ok {
		return token, fmt.Errorf("No token found")
	}
	return token, nil
}

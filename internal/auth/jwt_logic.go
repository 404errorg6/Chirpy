package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	expiry := jwt.NewNumericDate(time.Now().Add(expiresIn))

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  now,
		ExpiresAt: expiry,
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var id uuid.UUID
	var claims *jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return id, fmt.Errorf("Token parse error: %v\n", err)
	}
	if !token.Valid {
		return id, fmt.Errorf("Invalid token\n")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return id, fmt.Errorf("Unknown claims type\n")
	}

	idStr := claims.Subject
	id, err = uuid.Parse(idStr)
	if err != nil {
		return id, fmt.Errorf("Error converting to uuid: %v\n", err)
	}
	return id, nil
}

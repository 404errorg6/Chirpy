package auth

import (
	"testing"
)

func TestHash(t *testing.T) {
	password := "sup3r_s3cr3tP@a55word"
	hash, err := HashPassword(password)
	match, err := CheckPasswordHash(password, hash)
	if !match || err != nil {
		t.Errorf("\nHash: %v\nError: %v\n", hash, err)
	}
}

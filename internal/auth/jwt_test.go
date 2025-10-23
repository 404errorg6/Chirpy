package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var secretKey = "G0_t0_h311_m5"

func TestJWT(t *testing.T) {
	userID, err := uuid.Parse("652562d6-4244-4462-b079-151355c3846b")
	if err != nil {
		t.Fatal("UID parsing failed miserably\n")
	}

	signed, err := MakeJWT(userID, string(secretKey), time.Minute*10)
	if err != nil {
		t.Errorf("Error while creating JWT: %v\n", err)
	}

	gotID, err := ValidateJWT(signed, string(secretKey))
	if err != nil {
		t.Errorf("Validate error: %v\n", err)
	}
	if userID != gotID {
		t.Errorf("Invalid ID returned\nExpected: %v\tActual: %v\n", userID, gotID)
	}
}

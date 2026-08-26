package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
	"net/http"
)

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "test-secret"

	token, err := MakeJWT(id, tokenSecret, time.Minute)
	if err != nil {
		t.Fatalf("Token creation error: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Token validation error: %v", err)
	}

	if parsedID != id {
		t.Errorf("IDs differ: got %v, want %v", parsedID, id)
	}
}

func TestExpiredJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "test-secret"

	token, err := MakeJWT(id, tokenSecret, time.Second)
	if err != nil {
		t.Fatalf("Token creation error: %v", err)
	}

	time.Sleep(2 * time.Second)

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("Token validation error: expected err, got nil")
	}

	if parsedID == id {
		t.Errorf("IDs are the same: expected err")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "test-secret"
	wrongSecret := "wrong-secret"

	token, err := MakeJWT(id, wrongSecret, time.Minute)
	if err != nil {
		t.Fatalf("Token creation error: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("Token validation error: expected err, got nil")
	}

	if parsedID == id {
		t.Errorf("IDs are the same: expected err")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := make(http.Header)
	header.Set("Authorization", "Bearer testJwtToken")

	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("Bearer Token error: %v", err)
	}

	if token != "testJwtToken" {
		t.Fatalf("Token extraction error: got %v, expected testJwtToken", token)
	}
}


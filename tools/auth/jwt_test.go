package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := createJWT(jwt.MapClaims{
		"type": "test",
		"id": 12,
	}, secret, 3600)
	
	if err != nil {
		t.Errorf("error creating jwt: %v", err)
	}
	if token == "" {
		t.Errorf("expected token to not be empty")
	}
}
package auth

import "testing"

func TestHashPassword(t * testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Errorf("expected hash to not be empty")
	}

	if hash == "password" {
		t.Errorf("expected hash to be different to password")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !CheckPassword(hash, "password") {
		t.Errorf("expected password to match hash")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Errorf("expected wrong password to not match hash")
	}
}
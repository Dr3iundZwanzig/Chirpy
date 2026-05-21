package auth

import (
	"testing"
)

func TestEmptyHashedPassword(t *testing.T) {
	password := ""
	_, err := HashedPassword(password)
	if err == nil {
		t.Errorf("Error is nil want: password empty")
	}
}

func TestCheckPasswordHashMatch(t *testing.T) {
	password := "password"
	HashedPassword, err := HashedPassword(password)
	if err != nil {
		t.Errorf("Error creating password hash: %v", err)
	}
	ok, err := CheckPasswordHash(password, HashedPassword)
	if err != nil {
		t.Errorf("Error comparing password hash: %v", err)
	}
	if !ok {
		t.Errorf("Error password and has does not match")
	}
}

func TestCheckPasswordHashNoMatch(t *testing.T) {
	password := "password"
	HashedPassword, err := HashedPassword(password)
	if err != nil {
		t.Errorf("Error creating password hash: %v", err)
	}
	wrongPassword := "wrong"
	ok, err := CheckPasswordHash(wrongPassword, HashedPassword)
	if err == nil {
		t.Errorf("No error comparing password hash Wants: wrong password")
	}
	if ok {
		t.Errorf("Error password and has does match but it should not")
	}
}

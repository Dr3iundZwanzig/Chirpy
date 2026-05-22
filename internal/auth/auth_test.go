package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	accessToken := "accessToken"
	bearer := "Bearer " + accessToken
	req, err := http.NewRequest("GET", "local", nil)
	if err != nil {
		t.Errorf("internal error making new request: %v", err)
	}
	req.Header.Add("Authorization", bearer)
	token, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}
	if token != accessToken {
		t.Errorf("Wrong token: %v", err)
	}
}

func TestWrongFromatingGetBearerToken(t *testing.T) {
	accessToken := "accessToken"
	bearer := "Ber" + accessToken
	req, err := http.NewRequest("GET", "local", nil)
	if err != nil {
		t.Errorf("internal error making new request: %v", err)
	}
	req.Header.Add("Authorization", bearer)
	token, err := GetBearerToken(req.Header)
	if err == nil {
		t.Errorf("Error getting token: %v", err)
	}
	if token == accessToken {
		t.Errorf("Token corret should be an empty string: %v", err)
	}
}

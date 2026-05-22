package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashedPassword(password string) (string, error) {
	if len(password) == 0 {
		return password, fmt.Errorf("password empty")
	}
	hashPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return password, fmt.Errorf("Error creating password hash: %v", err)
	}
	return hashPassword, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password and hash: %v", err)
	}
	if !ok {
		return false, fmt.Errorf("wrong password")
	}
	return ok, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	TokenString := headers.Get("Authorization")
	if TokenString == "" {
		return "", fmt.Errorf("Error no token found in header")
	}
	tokenSplit := strings.Split(TokenString, " ")
	if len(tokenSplit) < 2 || tokenSplit[0] != "Bearer" {
		return "", fmt.Errorf("Error invalid token found in header")
	}
	return tokenSplit[1], nil
}

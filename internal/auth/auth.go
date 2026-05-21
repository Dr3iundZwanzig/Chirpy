package auth

import (
	"fmt"

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

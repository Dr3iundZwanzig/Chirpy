package auth

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret112"
	expiresIn, err := time.ParseDuration("10m")
	if err != nil {
		log.Fatalf("Internal error could not parse duration")
	}
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Error creating token string in MakeJWT function: %v", err)
	}
	validatedUserId, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("Error validating token in ValidateJWT function: %v", err)
	}
	if validatedUserId != userID {
		t.Errorf("User IDs not the same: %v", err)
	}
}

func TestTimeRunOutValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret112"
	expiresIn, err := time.ParseDuration("1ms")
	if err != nil {
		log.Fatalf("Internal error could not parse duration")
	}
	time.Sleep(2 * time.Second)
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Error creating token string in MakeJWT function: %v", err)
	}
	validatedUserId, err := ValidateJWT(tokenString, tokenSecret)
	if validatedUserId == userID {
		t.Errorf("User IDs is the same but should not be: %v", err)
	}
	if err == nil {
		t.Errorf("Error valdidated token true Expected to be false time ran out: %v", err)
	}
}

func TestWrongSecretValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret112"
	wrongSecret := "22sf"
	expiresIn, err := time.ParseDuration("10m")
	if err != nil {
		log.Fatalf("Internal error could not parse duration")
	}
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Error creating token string in MakeJWT function: %v", err)
	}
	validatedUserId, err := ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Errorf("Error validating token works with wrong secret key: %v", err)
	}
	if validatedUserId == userID {
		t.Errorf("Error validating token works with wrong secret key: %v", err)
	}
}

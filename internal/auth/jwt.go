package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenAccess string = "chirpy-access"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    TokenAccess,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedToken, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing token: %v", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	user := uuid.UUID{}
	type CustomClaim struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error validating token: %v", err)
	} else if claims, ok := token.Claims.(*CustomClaim); ok {
		user, err = uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("Error getting user ID: %v", err)
		}
		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			return uuid.UUID{}, fmt.Errorf("Error token expired")
		}
		if claims.Issuer != TokenAccess {
			return uuid.UUID{}, fmt.Errorf("Wrong issuer")
		}
	} else {
		return uuid.UUID{}, fmt.Errorf("Error getting claim: %v", err)
	}
	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("Error invalid token: %v", err)
	}
	return user, nil
}

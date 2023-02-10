package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(secret, id string, ttl time.Duration) (string, error) {
	// Define claims
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(ttl).Unix(),
		"iss": "rest-api",
	}

	// Create a new token with signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return token
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString, secret string) (*jwt.Token, jwt.MapClaims, error) {
	// Verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify token is signed with method HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return token, claims, nil
}

package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	id := "123"
	secret := "secret"
	token, err := CreateToken(secret, id, 5*time.Second)
	assert.NoError(t, err)
	jwtToken, claims, err := VerifyToken(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, jwtToken)
	assert.Equal(t, id, claims["id"])
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		label        string
		ttl          time.Duration
		createSecret string
		verifySecret string
		expectedErr  error
	}{
		{
			label:        "token is expired",
			ttl:          0 * time.Second,
			createSecret: "secret",
			verifySecret: "secret",
			expectedErr:  jwt.ErrTokenExpired,
		},
		{
			label:        "token is signed with another secret",
			ttl:          5 * time.Second,
			createSecret: "wrong",
			verifySecret: "secret",
			expectedErr:  jwt.ErrSignatureInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			token, err := CreateToken(test.createSecret, "123", test.ttl)
			assert.NoError(t, err)
			_, _, err = VerifyToken(token, test.verifySecret)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

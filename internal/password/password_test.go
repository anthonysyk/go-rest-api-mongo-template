package password

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	err = VerifyPassword(hashed, password)
	assert.NoError(t, err)
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		label       string
		password    string
		input       string
		expectedErr error
	}{
		{
			label:       "password is correct",
			password:    "secret",
			input:       "secret",
			expectedErr: nil,
		},
		{
			label:       "password is wrong",
			password:    "secret",
			input:       "wrong",
			expectedErr: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			hashed, err := HashPassword(test.password)
			assert.NoError(t, err)
			err = VerifyPassword(hashed, test.input)
			assert.ErrorIs(t, test.expectedErr, err)
		})
	}
}

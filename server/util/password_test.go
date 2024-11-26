package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPasword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPasword)

	err = CheckPassword(password, hashedPasword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPasword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

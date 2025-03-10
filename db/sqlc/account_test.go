package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccout(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  10,
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
}

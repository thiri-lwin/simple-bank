package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/thiri-bank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.Owner, arg.Owner)
	require.Equal(t, acc.Currency, arg.Currency)
	require.Equal(t, acc.Balance, arg.Balance)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAcc := createRandomAccount(t)

	acc, err := testQueries.GetAccount(context.Background(), createdAcc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, createdAcc.ID, acc.ID)
	require.Equal(t, createdAcc.Owner, acc.Owner)
	require.Equal(t, createdAcc.Currency, acc.Currency)
	require.WithinDuration(t, createdAcc.CreatedAt, acc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomBalance(),
	}

	err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, account.ID, acc.ID)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, account.Owner, acc.Owner)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}

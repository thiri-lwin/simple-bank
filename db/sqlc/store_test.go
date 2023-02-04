package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransaferTx(ctx, TransaferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		dbTransfer, err := store.GetTransaction(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, dbTransfer)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, amount, fromEntry.Amount)
		require.Equal(t, DEBIT, fromEntry.TransactionType)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.Equal(t, CREDIT, toEntry.TransactionType)

		// check accounts
		fromAcc := result.FromAccount
		require.NotEmpty(t, fromAcc)
		require.Equal(t, account1.ID, fromAcc.ID)
		require.Equal(t, account1.Owner, fromAcc.Owner)
		//require.Equal(t, account1.Balance-amount, fromAcc.Balance)

		toAcc := result.ToAccount
		require.NotEmpty(t, toAcc)
		require.Equal(t, account2.ID, toAcc.ID)
		require.Equal(t, account2.Owner, toAcc.Owner)
		//require.Equal(t, account2.Balance+amount, toAcc.Balance)

		// check balances
		diff1 := account1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
	}

	updatedAcc, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, updatedAcc.Balance, account1.Balance-(int64(n)*amount))

	updatedAcc2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc2)
	require.Equal(t, updatedAcc2.Balance, account2.Balance+(int64(n)*amount))
}

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/thiri-bank/util"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func createRandomEntry(t *testing.T) Entry {
	acc := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID:       acc.ID,
		Amount:          util.RandomBalance(),
		TransactionType: util.RandomTransactionType(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.TransactionType, entry.TransactionType)
	require.NotZero(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)
	return entry
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	dbRes, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbRes)
	require.Equal(t, entry.AccountID, dbRes.AccountID)
	require.Equal(t, entry.Amount, dbRes.Amount)
	require.Equal(t, entry.TransactionType, dbRes.TransactionType)
	require.Equal(t, entry.ID, dbRes.ID)
	require.WithinDuration(t, entry.CreatedAt, dbRes.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)
	for _, acc := range entries {
		require.NotEmpty(t, acc)
	}
}

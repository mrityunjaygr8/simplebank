package db

import (
	"context"
	"testing"
	"time"

	"github.com/mrityunjaygr8/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.ID, entry1.ID)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.Amount, entry1.Amount)
	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)

}

func TestListEntry(t *testing.T) {
	for x := 0; x < 10; x++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, len(entries), 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListEntryForAccount(t *testing.T) {
	account := createRandomAccount(t)
	entries := []Entry{}

	for x := 0; x < 10; x++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    utils.RandomMoney(),
		}

		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)

		entries = append(entries, entry)
	}

	arg := ListEntriesForAccountParams{
		Limit:     5,
		Offset:    5,
		AccountID: account.ID,
	}
	entries, err := testQueries.ListEntriesForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, len(entries), 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

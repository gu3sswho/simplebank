package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gu3sswho/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entryActual := createRandomEntry(t)

	entryExpect, err := testQueries.GetEntry(context.Background(), entryActual.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryExpect)

	require.Equal(t, entryExpect.ID, entryActual.ID)
	require.Equal(t, entryExpect.AccountID, entryActual.AccountID)
	require.Equal(t, entryExpect.Amount, entryActual.Amount)
	require.WithinDuration(t, entryExpect.CreatedAt, entryActual.CreatedAt, time.Second)
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
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	entryBefore := createRandomEntry(t)

	updateArg := UpdateEntryParams{
		ID:     entryBefore.ID,
		Amount: util.RandomMoney(),
	}

	entryAfter, err := testQueries.UpdateEntry(context.Background(), updateArg)

	require.NoError(t, err)
	require.NotEmpty(t, entryAfter)

	require.Equal(t, entryBefore.ID, entryAfter.ID)
	require.Equal(t, entryBefore.AccountID, entryAfter.AccountID)
	require.Equal(t, updateArg.Amount, entryAfter.Amount)
	require.WithinDuration(t, entryBefore.CreatedAt, entryAfter.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	deletedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedEntry)
}

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gu3sswho/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transferActual := createRandomTransfer(t)

	transferExpect, err := testQueries.GetTransfer(context.Background(), transferActual.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferExpect)

	require.Equal(t, transferExpect.ID, transferActual.ID)
	require.Equal(t, transferExpect.FromAccountID, transferActual.FromAccountID)
	require.Equal(t, transferExpect.ToAccountID, transferActual.ToAccountID)
	require.Equal(t, transferExpect.Amount, transferActual.Amount)
	require.WithinDuration(t, transferExpect.CreatedAt, transferActual.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transferBefore := createRandomTransfer(t)

	updateArg := UpdateTransferParams{
		ID:     transferBefore.ID,
		Amount: util.RandomMoney(),
	}

	transferAfter, err := testQueries.UpdateTransfer(context.Background(), updateArg)

	require.NoError(t, err)
	require.NotEmpty(t, transferAfter)

	require.Equal(t, transferBefore.ID, transferAfter.ID)
	require.Equal(t, transferBefore.FromAccountID, transferAfter.FromAccountID)
	require.Equal(t, transferBefore.ToAccountID, transferAfter.ToAccountID)
	require.Equal(t, updateArg.Amount, transferAfter.Amount)
	require.WithinDuration(t, transferBefore.CreatedAt, transferAfter.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	deletedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedTransfer)
}

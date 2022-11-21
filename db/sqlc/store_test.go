package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	//run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)

		go func() {
			result, err := store.TransferTx(ctx, TransferTxParam{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errChan <- err
			resultChan <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		from := result.FromAccount

		require.NotEmpty(t, from)
		require.Equal(t, fromAccount.ID, from.ID)
		require.Equal(t, fromAccount.Owner, from.Owner)
		require.Equal(t, fromAccount.Currency, from.Currency)
		require.WithinDuration(t, fromAccount.CreatedAt, from.CreatedAt, time.Second)

		_, err = store.GetAccount(context.Background(), from.ID)
		require.NoError(t, err)

		to := result.ToAccount

		require.NotEmpty(t, to)
		require.Equal(t, toAccount.ID, to.ID)
		require.Equal(t, toAccount.Owner, to.Owner)
		require.Equal(t, toAccount.Currency, to.Currency)
		require.WithinDuration(t, toAccount.CreatedAt, to.CreatedAt, time.Second)

		_, err = store.GetAccount(context.Background(), to.ID)
		require.NoError(t, err)

		//check balance
		diff1 := fromAccount.Balance - from.Balance
		diff2 := to.Balance - toAccount.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)         //always positive balance
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount ...

		k := int(diff1 / amount)
		fmt.Printf("diff1=%d, amount=%d, k=%d\n", diff1, amount, k)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	//run n concurrent transfer transactions (n - should be even)
	n := 10
	amount := int64(10)

	errChan := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := fromAccount.ID
		toAccountID := toAccount.ID

		if i%2 == 1 {
			fromAccountID = toAccount.ID
			toAccountID = fromAccount.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errChan <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
	}

	//check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance, updatedAccount2.Balance)
}

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gu3sswho/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountActual := createRandomAccount(t)

	accountExpect, err := testQueries.GetAccount(context.Background(), accountActual.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountExpect)

	require.Equal(t, accountExpect.ID, accountActual.ID)
	require.Equal(t, accountExpect.Owner, accountActual.Owner)
	require.Equal(t, accountExpect.Balance, accountActual.Balance)
	require.Equal(t, accountExpect.Currency, accountActual.Currency)
	require.WithinDuration(t, accountExpect.CreatedAt, accountActual.CreatedAt, time.Second)
}

func TestGetAccountForUpdate(t *testing.T) {
	accountActual := createRandomAccount(t)

	accountExpect, err := testQueries.GetAccountForUpdate(context.Background(), accountActual.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountExpect)

	require.Equal(t, accountExpect.ID, accountActual.ID)
	require.Equal(t, accountExpect.Owner, accountActual.Owner)
	require.Equal(t, accountExpect.Balance, accountActual.Balance)
	require.Equal(t, accountExpect.Currency, accountActual.Currency)
	require.WithinDuration(t, accountExpect.CreatedAt, accountActual.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestUpdateAccount(t *testing.T) {
	accountBefore := createRandomAccount(t)

	updateArg := UpdateAccountParams{
		ID:      accountBefore.ID,
		Balance: util.RandomMoney(),
	}

	accountAfter, err := testQueries.UpdateAccount(context.Background(), updateArg)

	require.NoError(t, err)
	require.NotEmpty(t, accountAfter)

	require.Equal(t, accountBefore.ID, accountAfter.ID)
	require.Equal(t, accountBefore.Owner, accountAfter.Owner)
	require.Equal(t, updateArg.Balance, accountAfter.Balance)
	require.Equal(t, accountBefore.Currency, accountAfter.Currency)
	require.WithinDuration(t, accountBefore.CreatedAt, accountAfter.CreatedAt, time.Second)
}

func TestAddAccountBalance(t *testing.T) {
	accountBefore := createRandomAccount(t)

	addArg := AddAccountBalanceParams{
		ID:     accountBefore.ID,
		Amount: util.RandomMoney(),
	}

	accountAfter, err := testQueries.AddAccountBalance(context.Background(), addArg)

	require.NoError(t, err)
	require.NotEmpty(t, accountAfter)

	require.Equal(t, accountBefore.ID, accountAfter.ID)
	require.Equal(t, accountBefore.Owner, accountAfter.Owner)
	require.Equal(t, accountBefore.Balance+addArg.Amount, accountAfter.Balance)
	require.Equal(t, accountBefore.Currency, accountAfter.Currency)
	require.WithinDuration(t, accountBefore.CreatedAt, accountAfter.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

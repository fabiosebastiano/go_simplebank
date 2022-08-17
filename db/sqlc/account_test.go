package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fabiosebastiano/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username, //randomly generated
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, error := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, error)
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
	account := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.EqualValues(t, account.ID, account2.ID)
	require.EqualValues(t, account.Balance, account2.Balance)
	require.EqualValues(t, account.Currency, account2.Currency)
	require.WithinDuration(t, account.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	params := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomInt(1, 1000),
	}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.NotEqual(t, account.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

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
		require.Equal(t, account.Owner, lastAccount.Owner)
	}
}

package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrecy(),
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
	accountA := createRandomAccount(t)
	accountB, err := testQueries.GetAccount(context.Background(), accountA.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountB)
	require.Equal(t, accountA.ID, accountB.ID)
	require.Equal(t, accountA.Balance, accountB.Balance)
	require.Equal(t, accountA.Owner, accountB.Owner)
	require.WithinDuration(t, accountA.CreatedAt, accountB.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	accountA := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      accountA.ID,
		Balance: util.RandomMoney(),
	}

	accountB, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accountB)

	require.Equal(t, accountA.ID, accountB.ID)
	require.Equal(t, accountA.Owner, accountB.Owner)
	require.Equal(t, accountA.Currency, accountB.Currency)
	require.Equal(t, arg.Balance, accountB.Balance)
	require.WithinDuration(t, accountA.CreatedAt, accountB.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accountA := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	accountB, err := testQueries.GetAccount(context.Background(), accountA.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountB)
}

func TestListAccounts(t *testing.T) {
	for a := 0; a < 10; a++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)
}

package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID) 

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t) 
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t) 
	transfer1 := createRandomTransfer(t, account1, account2)
	
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, transfer2) 
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
}


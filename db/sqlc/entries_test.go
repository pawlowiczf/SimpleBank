package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID) 
	require.NotZero(t, entry.CreatedAt) 

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t) 
	entryA := createRandomEntry(t, account) 
	entryB, err := testQueries.GetEntry(context.Background(), entryA.ID)

	require.NoError(t, err) 
	require.NotEmpty(t, entryB) 

	require.Equal(t, entryA.ID, entryB.ID) 
	require.Equal(t, entryA.AccountID, entryB.AccountID) 
	require.Equal(t, entryA.Amount, entryB.Amount) 
}
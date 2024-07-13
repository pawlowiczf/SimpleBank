package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for a := 0; a < n; a++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for a := 0; a < n; a++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, accountA.ID)
		require.Equal(t, transfer.ToAccountID, accountB.ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries:
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountA.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountB.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountA.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountB.ID, toAccount.ID)

		// check accounts' balance
		diff1 := accountA.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accountB.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccountA, err := testQueries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testQueries.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)

	require.Equal(t, accountA.Balance-int64(n)*amount, updatedAccountA.Balance)
	require.Equal(t, accountB.Balance+int64(n)*amount, updatedAccountB.Balance)
}

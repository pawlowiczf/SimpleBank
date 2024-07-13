package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries // embedding
	db       *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// execTx executes a function withing a database transaction
func (store *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, adds account entries, and updates accounts' balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult 

	err := store.execTx(ctx, func(q *Queries) error {
		var err error 

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg)) // tutaj ciekawa konwersja. Mozna by oczywiscie zrobić CreateTransferParams{}, ale ta struktura ma takie same pola.
		if err != nil {
			return nil 
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID, 
			Amount: -arg.Amount,
		})
		if err != nil {
			return nil 
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID, 
			Amount: arg.Amount,
		})
		if err != nil {
			return nil 
		}

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err 
		}


		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err 
		}


		// TODO: update accounts' balance
		return nil
	})

	return result, err 
}

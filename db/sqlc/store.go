package db

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	CREDIT = "Credit"
	DEBIT  = "Debit"
)

var txKey = struct{}{}

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

type TransaferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferTxResult struct {
	Transfer    Transaction
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (s *Store) TransaferTx(ctx context.Context, arg TransaferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transaction")
		result.Transfer, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:       arg.FromAccountID,
			Amount:          arg.Amount,
			TransactionType: DEBIT,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:       arg.ToAccountID,
			Amount:          arg.Amount,
			TransactionType: CREDIT,
		})

		fmt.Println(txName, "get account 1")
		fromAccount, err := q.GetAccount(context.Background(), arg.FromAccountID)
		if err != nil {
			return err
		}

		updateFromAccArg := UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: fromAccount.Balance - arg.Amount,
		}
		fmt.Println(txName, "update account 1")
		err = q.UpdateAccount(context.Background(), updateFromAccArg)
		if err != nil {
			return err
		}
		result.FromAccount = Account{
			ID:        arg.FromAccountID,
			Owner:     fromAccount.Owner,
			Balance:   updateFromAccArg.Balance,
			Currency:  fromAccount.Currency,
			CreatedAt: fromAccount.CreatedAt,
		}

		fmt.Println(txName, "get account 2")
		toAccount, err := q.GetAccount(context.Background(), arg.ToAccountID)
		if err != nil {
			return err
		}

		updateToAccArg := UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: toAccount.Balance + arg.Amount,
		}
		fmt.Println(txName, "update account 2")
		err = q.UpdateAccount(context.Background(), updateToAccArg)
		if err != nil {
			return err
		}
		result.ToAccount = Account{
			ID:        arg.ToAccountID,
			Owner:     toAccount.Owner,
			Balance:   updateToAccArg.Balance,
			Currency:  toAccount.Currency,
			CreatedAt: toAccount.CreatedAt,
		}

		return nil
	})

	return result, err
}

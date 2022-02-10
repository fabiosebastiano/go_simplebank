package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store contiene tutte le funzionalità per eseguire query e transazioni
type Store struct {
	*Queries
	db *sql.DB
}

//NewStore crea un nuovo oggetto Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

var txKey = struct{}{}

//execTx chiama la fn di callback sulla transazione appena creata - privata perchè non si usi fuori dal pkg
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//INIZIARE TX
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//CREA QUERY ALL'INTERNO DELLA TX
	query := New(tx)
	err = fn(query)
	if err != nil {
		// errore -> rollback
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackError)
			//return fmt.Errorf("Tx error is %w & rollback error is %w", err, rollbackError)
		}
		return err
	}

	//ESEGUE COMMIT
	return tx.Commit()
}

//TransferTxParams Contiene tutti i parametri necessari a completare la transazione
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult contiene il risultato della transazione
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx performa un trasferimento di denaro tra due conti
/*
	- crea un record TRANSFER
	- crea due record ENTRY
	- modifica il balance dei due account
*/
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)

		} else {

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)

		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	queries *Queries,
	fromAccountID int64,
	toAccountID int64,
	fromAmount int64,
	toAmount int64) (
	fromAccount Account, toAccount Account, err error) {

	fromAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}
	toAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})

	return
}

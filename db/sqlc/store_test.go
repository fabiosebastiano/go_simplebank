package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	//fmt.Println("--> BEFORE", accountFrom.Balance, accountTo.Balance)

	// per verificare che le TRANSAZIONI funzionino correttamente, si usano goroutine concorrenti
	n := 5
	amount := int64(10)

	//per testare che le goroutine non vadano in errore, utilizzo i CHAN per avere indietro gli eventuali errori
	errorsChan := make(chan error)
	resultsChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})
			errorsChan <- err
			resultsChan <- result
		}()
	}

	// VERIFICARE SE E' ARRIVATA ROBA NEI CANALI
	for j := 0; j < n; j++ {
		//controllo che non ci siano errori
		err := <-errorsChan
		require.NoError(t, err)

		//controllo che ci siano risultati
		res := <-resultsChan
		require.NotEmpty(t, res)

		//controllo che la transazione sia effettivamente stata creata
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//verifico che sia avvenuta davvero la scrittura a db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//verifico le entries
		entryFrom := res.FromEntry
		require.NotEmpty(t, entryFrom)
		require.Equal(t, accountFrom.ID, entryFrom.AccountID)
		require.EqualValues(t, -amount, entryFrom.Amount)
		require.NotZero(t, entryFrom.ID)
		require.NotZero(t, entryFrom.CreatedAt)
		//verifico che sia avvenuta davvero la scrittura a db
		_, err = store.GetEntry(context.Background(), entryFrom.ID)
		require.NoError(t, err)

		entryTo := res.ToEntry
		require.NotEmpty(t, entryTo)
		require.Equal(t, accountTo.ID, entryTo.AccountID)
		require.EqualValues(t, amount, entryTo.Amount)
		require.NotZero(t, entryTo.ID)
		require.NotZero(t, entryTo.CreatedAt)

		//verifico che sia avvenuta davvero la scrittura a db
		_, err = store.GetEntry(context.Background(), entryTo.ID)
		require.NoError(t, err)

		// TODO  check account's balances
		accountFromTx := res.FromAccount
		require.NotEmpty(t, accountFromTx)
		require.Equal(t, accountFrom.ID, accountFromTx.ID)

		accountToTx := res.ToAccount
		require.NotEmpty(t, accountToTx)
		require.Equal(t, accountTo.ID, accountToTx.ID)

		//fmt.Println("--> DURING", accountFromTx.Balance, accountToTx.Balance)

		diff1 := accountFrom.Balance - accountFromTx.Balance
		diff2 := accountToTx.Balance - accountTo.Balance
		require.EqualValues(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

	}

	updatedAccountFrom, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountFrom)
	require.EqualValues(t, updatedAccountFrom.Balance, accountFrom.Balance-int64(n)*amount)

	updatedAccountTo, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountTo)
	require.EqualValues(t, updatedAccountTo.Balance, accountTo.Balance+int64(n)*amount)

	//fmt.Println("--> AFTER", updatedAccountFrom.Balance, updatedAccountTo.Balance, int64(n)*amount)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	//fmt.Println("--> BEFORE", accountFrom.Balance, accountTo.Balance)

	// per verificare che le TRANSAZIONI funzionino correttamente, si usano goroutine concorrenti
	n := 10
	amount := int64(10)

	//per testare che le goroutine non vadano in errore, utilizzo i CHAN per avere indietro gli eventuali errori
	errorsChan := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountFrom.ID
		toAccountID := accountTo.ID

		if i%2 == 1 {
			fromAccountID = accountTo.ID
			toAccountID = accountFrom.ID
		}

		go func() {

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errorsChan <- err
		}()
	}

	// VERIFICARE SE E' ARRIVATA ROBA NEI CANALI
	for j := 0; j < n; j++ {
		//controllo che non ci siano errori
		err := <-errorsChan
		require.NoError(t, err)

	}

	updatedAccountFrom, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountFrom)

	updatedAccountTo, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountTo)

	require.EqualValues(t, updatedAccountFrom.Balance, accountFrom.Balance)
	require.EqualValues(t, updatedAccountTo.Balance, accountTo.Balance)
}

package db

import (
	"context"
	"testing"
	"time"

	"github.com/fabiosebastiano/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: 1, //randomly generated
		Amount:    util.RandomMoney(),
	}

	entry, error := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, error)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

func TestCreateEntry(t *testing.T) {
	CreateEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := CreateEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.EqualValues(t, entry.ID, entry2.ID)
	require.EqualValues(t, entry.AccountID, entry2.AccountID)

	require.EqualValues(t, entry.Amount, entry2.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateEntry(t)
	}
	param := ListEntriesParams{
		AccountID: 1,
		Limit:     5,
		Offset:    5,
	}
	entrys, err := testQueries.ListEntries(context.Background(), param)
	require.NoError(t, err)
	require.NotEmpty(t, entrys)
}

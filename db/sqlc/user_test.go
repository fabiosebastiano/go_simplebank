package db

import (
	"context"
	"testing"
	"time"

	"github.com/fabiosebastiano/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomUsername(), //randomly generated
		HashedPassword: "secret",
		FullName:       util.RandomFullname(),
		Email:          util.RandomEmail(),
	}

	user, error := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, error)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.EqualValues(t, user.Username, user2.Username)
	require.EqualValues(t, user.FullName, user2.FullName)
	require.EqualValues(t, user.Email, user2.Email)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
}

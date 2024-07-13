package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err) 

	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		HashedPassword:  hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	userA := createRandomUser(t)
	userB, err := testQueries.GetUser(context.Background(), userA.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userB)
	require.Equal(t, userA.Username, userB.Username)
	require.Equal(t, userA.Email, userB.Email)
	require.Equal(t, userA.HashedPassword, userB.HashedPassword)
	require.WithinDuration(t, userA.CreatedAt, userB.CreatedAt, time.Second)
}

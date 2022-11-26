package db

import (
	"context"
	"testing"
	"time"

	"github.com/gu3sswho/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(15))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	userActual := createRandomUser(t)

	userExpect, err := testQueries.GetUser(context.Background(), userActual.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userExpect)

	require.Equal(t, userExpect.Username, userActual.Username)
	require.Equal(t, userExpect.HashedPassword, userActual.HashedPassword)
	require.Equal(t, userExpect.FullName, userActual.FullName)
	require.Equal(t, userExpect.Email, userActual.Email)
	require.WithinDuration(t, userExpect.CreatedAt, userActual.CreatedAt, time.Second)
	require.WithinDuration(t, userExpect.PasswordChangedAt, userActual.PasswordChangedAt, time.Second)
}

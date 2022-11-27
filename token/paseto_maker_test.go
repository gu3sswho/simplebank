package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/gu3sswho/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestNewPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.ID)
}

func TestExpiredPaseto(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPaseto(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken(util.RandomString(32))
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestTooShortKeyLength(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(16))

	require.Error(t, err)
	require.EqualError(t, err, fmt.Errorf("invalid secret key length: must be at least %d characters", chacha20poly1305.KeySize).Error())
	require.Nil(t, maker)
}

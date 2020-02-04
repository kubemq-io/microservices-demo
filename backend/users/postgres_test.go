package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgres_Functions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	name := "lior.nabat@gmail.com"
	password := "some_password"
	email := "lior.nabat@gmail.com"
	cfg, err := LoadConfig()
	require.NoError(t, err)
	pq, err := NewPostgres(cfg)
	require.NoError(t, err)
	u, err := pq.Register(ctx, &NewUser{
		Name:     name,
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	require.Equal(t, u.Name, name)
	require.Equal(t, u.Email, email)
	require.Equal(t, u.Password, "")
	require.Equal(t, u.Id, "")
	require.Equal(t, u.State, 1)
	require.NotEmpty(t, u.Token)

	_, err = pq.Register(ctx, &NewUser{
		Name:     name,
		Password: password,
		Email:    email,
	})
	require.Error(t, err)

	_, err = pq.Login(ctx, &LoginRequest{
		Name:     name,
		Password: password,
	})
	require.Error(t, err)
	err = pq.VerifyRegistration(ctx, &VerifyRegistration{
		Name:  u.Name,
		Token: u.Token,
	})
	require.NoError(t, err)
	loginResp, err := pq.Login(ctx, &LoginRequest{
		Name:     name,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Token)
	require.NotEmpty(t, loginResp.Token)
	require.NotZero(t, loginResp.Expiry.Sub(time.Now()))

	logoutResp, err := pq.Logout(ctx, &LogoutRequest{
		UserId: loginResp.UserId,
	})
	require.NoError(t, err)
	require.Equal(t, loginResp.Token, logoutResp.Token)
	require.NotEmpty(t, logoutResp.UserId)

	resetResp, err := pq.PasswordResetRequest(ctx, &PasswordResetRequest{
		Name: name,
	})
	require.NoError(t, err)
	require.Equal(t, resetResp.Email, email)
	require.Equal(t, resetResp.UserID, loginResp.UserId)
	require.Equal(t, resetResp.Name, name)
	require.NotEmpty(t, resetResp.Token)

	_, err = pq.Login(ctx, &LoginRequest{
		Name:     name,
		Password: password,
	})
	require.Error(t, err)
	changeResp, err := pq.PasswordChangeRequest(ctx, &PasswordChangeRequest{
		Name:        name,
		Token:       resetResp.Token,
		NewPassword: "new_password",
	})
	require.NoError(t, err)
	require.Equal(t, changeResp.Name, name)
	require.Equal(t, changeResp.Email, email)
	require.Equal(t, changeResp.UserID, loginResp.UserId)
	newLogin, err := pq.Login(ctx, &LoginRequest{
		Name:     name,
		Password: "new_password",
	})
	require.NoError(t, err)
	require.NotEmpty(t, newLogin.UserId)
	require.NotEmpty(t, newLogin.Token)
	require.NotZero(t, newLogin.Expiry.Sub(time.Now()))

	lockResp, err := pq.Lock(ctx, &LockRequest{
		UserID: newLogin.UserId,
	})
	require.NoError(t, err)
	require.Equal(t, lockResp.UserID, newLogin.UserId)
	require.Equal(t, lockResp.Email, email)
	require.Equal(t, lockResp.Name, name)

	_, err = pq.PasswordResetRequest(ctx, &PasswordResetRequest{
		Name: name,
	})
	require.Error(t, err)
	unlockResp, err := pq.Unlock(ctx, &UnlockRequest{
		UserID: newLogin.UserId,
	})
	require.NoError(t, err)
	require.Equal(t, unlockResp.UserID, newLogin.UserId)
	require.Equal(t, unlockResp.Email, email)
	require.Equal(t, unlockResp.Name, name)

}

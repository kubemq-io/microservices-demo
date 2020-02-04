package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"
)

func PrettyJson(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}

type User struct {
	Id        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	State     int       `json:"state" db:"state"`
	Token     string    `json:"token" db:"token"`
}

func (u *User) CheckState() error {
	switch u.State {
	case 1:
		return errors.New("user not verified")
	case 2:
		return nil
	case 3:
		return errors.New("user changed password")
	case 4:
		return errors.New("user locked")
	default:
		return errors.New("user invalid state")
	}

}
func (u *User) Data() []byte {
	data, _ := json.Marshal(u)
	return data
}

func getUser(data []byte) (*User, error) {
	u := &User{}
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type Login struct {
	UserId    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ExpireAt  time.Time `json:"expire_at" db:"expire_at"`
}
type NewUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func getNewUser(data []byte) (*NewUser, error) {
	nu := &NewUser{}
	err := json.Unmarshal(data, nu)
	if err != nil {
		return nil, err
	}
	return nu, nil
}

type VerifyRegistration struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func getVerificationRegistration(data []byte) (*VerifyRegistration, error) {
	vr := &VerifyRegistration{}
	err := json.Unmarshal(data, vr)
	if err != nil {
		return nil, err
	}
	return vr, nil

}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func getLoginRequest(data []byte) (*LoginRequest, error) {
	lr := &LoginRequest{}
	err := json.Unmarshal(data, lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}

type LoginResponse struct {
	UserId string    `json:"user_id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

func (lr *LoginResponse) Data() []byte {
	data, _ := json.Marshal(lr)
	return data
}

func getLoginResponse(data []byte) (*LoginResponse, error) {
	lr := &LoginResponse{}
	err := json.Unmarshal(data, lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}

type LogoutRequest struct {
	UserId string `json:"user_id"`
}

func getLogoutRequest(data []byte) (*LogoutRequest, error) {
	lr := &LogoutRequest{}
	err := json.Unmarshal(data, lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}

type LogoutResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

func (lr *LogoutResponse) Data() []byte {
	data, _ := json.Marshal(lr)
	return data
}

type PasswordResetRequest struct {
	Name string `json:"name"`
}

func getPasswordResetRequest(data []byte) (*PasswordResetRequest, error) {
	prr := &PasswordResetRequest{}
	err := json.Unmarshal(data, prr)
	if err != nil {
		return nil, err
	}
	return prr, nil
}

type PasswordResetResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (prr *PasswordResetResponse) Data() []byte {
	data, _ := json.Marshal(prr)
	return data
}

type PasswordChangeRequest struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func getPasswordChaneRequest(data []byte) (*PasswordChangeRequest, error) {
	pcr := &PasswordChangeRequest{}
	err := json.Unmarshal(data, pcr)
	if err != nil {
		return nil, err
	}
	return pcr, nil
}

type PasswordChangeResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (pcr *PasswordChangeResponse) Data() []byte {
	data, _ := json.Marshal(pcr)
	return data
}

type LockRequest struct {
	UserID string `json:"user_id"`
}

func getLockRequest(data []byte) (*LockRequest, error) {
	lr := &LockRequest{}
	err := json.Unmarshal(data, lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}

type LockResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (lr *LockResponse) Data() []byte {
	data, _ := json.Marshal(lr)
	return data
}

type UnlockRequest struct {
	UserID string `json:"user_id"`
}

func getUnlockRequest(data []byte) (*UnlockRequest, error) {
	ur := &UnlockRequest{}
	err := json.Unmarshal(data, ur)
	if err != nil {
		return nil, err
	}
	return ur, nil
}

type UnlockResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (ur *UnlockResponse) Data() []byte {
	data, _ := json.Marshal(ur)
	return data
}

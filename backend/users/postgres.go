package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var schema = `
DROP table IF EXISTS users;
DROP table IF EXISTS logins;
create table if not exists users
(
	-- Only integer types can be auto increment
	id varchar(255) not null
		constraint users_pk
			primary key,
	created_at timestamp not null,
	name varchar(255) not null,
	password varchar(255) not null,
	email varchar(255) not null,
	state integer default 0 not null,
	token varchar(255) not null
);

alter table users owner to postgres;

create unique index if not exists users_name_uindex
	on users (name);

create unique index if not exists users_id_uindex
	on users (id);

create table if not exists logins
(
	user_id varchar(255) not null
		constraint logins_pk
			primary key,
	token varchar(255) not null,
	updated_at timestamp not null,
	expire_at timestamp not null
);

alter table logins owner to postgres;

create unique index logins_user_id_uindex
	on logins (user_id);
`

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(config *Config) (*Postgres, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	log.Println(connString)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	db.MustExec(schema)

	p := &Postgres{
		db: db,
	}
	return p, nil
}

func (pq *Postgres) Register(ctx context.Context, newUser *NewUser) (*User, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE name=$1", newUser.Name)
	if err == nil {
		return nil, errors.New("user already exist")
	}
	hash, _ := HashPassword(newUser.Password)
	token := uuid.New().String()
	id := uuid.New().String()
	_, err = pq.db.ExecContext(ctx, "INSERT INTO users (id, created_at, name, password, email,state,token) VALUES ($1, $2, $3,$4,$5,$6,$7)", id, time.Now(), newUser.Name, hash, newUser.Email, 1, token)
	if err != nil {
		return nil, err
	}
	err = pq.db.Get(u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, errors.New("user created but not found")
	}
	u.Password = ""
	return u, nil
}

func (pq *Postgres) VerifyRegistration(ctx context.Context, vr *VerifyRegistration) error {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE name=$1 and token=$2", vr.Name, vr.Token)
	if err != nil {
		return errors.New("user not found or invalid token")
	}
	_, err = pq.db.ExecContext(ctx, "UPDATE  users SET state=2,token='' WHERE name=$1", vr.Name)
	if err != nil {
		return err
	}
	return nil
}

func (pq *Postgres) Login(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE name=$1", loginRequest.Name)
	if err != nil {
		return nil, errors.New("user not found or invalid password")
	}
	if err := u.CheckState(); err != nil {
		return nil, err
	}

	if !CheckPasswordHash(loginRequest.Password, u.Password) {
		return nil, errors.New("user not found or invalid password")
	}

	login := &Login{}

	exist := pq.db.Get(login, "SELECT * FROM logins WHERE user_id=$1", u.Id)
	if exist != nil {
		login = &Login{
			UserId:    u.Id,
			Token:     uuid.New().String(),
			UpdatedAt: time.Now(),
			ExpireAt:  time.Now().Add(30 * time.Minute),
		}
		_, err = pq.db.ExecContext(ctx, "INSERT INTO logins (user_id,token, updated_at, expire_at) VALUES ($1, $2, $3,$4)", login.UserId, login.Token, login.UpdatedAt, login.ExpireAt)
		if err != nil {
			return nil, err
		}
	} else {
		login.UpdatedAt = time.Now()
		login.ExpireAt = login.UpdatedAt.Add(30 * time.Minute)
		_, err = pq.db.ExecContext(ctx, "UPDATE  logins SET updated_at=$2,expire_at=$3 WHERE user_id=$1", login.UserId, login.UpdatedAt, login.ExpireAt)
		if err != nil {
			return nil, err
		}
	}
	lt := &LoginResponse{
		UserId: login.UserId,
		Expiry: login.ExpireAt,
	}
	return lt, nil
}

func (pq *Postgres) Logout(ctx context.Context, logoutRequest *LogoutRequest) (*LogoutResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE id=$1", logoutRequest.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := u.CheckState(); err != nil {
		return nil, err
	}

	login := &Login{}
	exist := pq.db.Get(login, "SELECT * FROM logins WHERE user_id=$1", logoutRequest.UserId)
	if exist == nil {
		_, err = pq.db.ExecContext(ctx, "DELETE FROM logins  WHERE user_id=$1", logoutRequest.UserId)
		return &LogoutResponse{
			Name:   u.Name,
			UserId: logoutRequest.UserId,
			Token:  login.Token,
		}, nil
	} else {
		return nil, errors.New("user not logged")
	}
}

func (pq *Postgres) PasswordResetRequest(ctx context.Context, req *PasswordResetRequest) (*PasswordResetResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE name=$1", req.Name)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := u.CheckState(); err != nil {
		return nil, err
	}
	token := uuid.New().String()
	_, err = pq.db.ExecContext(ctx, "UPDATE  users SET state=3 , token=$2 WHERE id=$1", u.Id, token)
	if err != nil {
		return nil, err
	}
	_, _ = pq.db.ExecContext(ctx, "DELETE FROM logins  WHERE user_id=$1", u.Id)
	return &PasswordResetResponse{
		UserID: u.Id,
		Name:   u.Name,
		Email:  u.Email,
		Token:  token,
	}, nil
}

func (pq *Postgres) PasswordChangeRequest(ctx context.Context, req *PasswordChangeRequest) (*PasswordChangeResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE name=$1", req.Name)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.State != 3 {
		return nil, errors.New("user didn't reset password")
	}
	if u.Token != req.Token {
		return nil, errors.New("invalid reset password token")
	}

	hash, _ := HashPassword(req.NewPassword)

	_, err = pq.db.ExecContext(ctx, "UPDATE  users SET password=$2, state=2 , token='' WHERE id=$1", u.Id, hash)
	if err != nil {
		return nil, err
	}

	return &PasswordChangeResponse{
		UserID: u.Id,
		Name:   u.Name,
		Email:  u.Email,
	}, nil
}

func (pq *Postgres) Lock(ctx context.Context, req *LockRequest) (*LockResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE id=$1", req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := u.CheckState(); err != nil {
		return nil, err
	}
	_, err = pq.db.ExecContext(ctx, "UPDATE  users SET  state=4  WHERE id=$1", u.Id)
	if err != nil {
		return nil, err
	}

	return &LockResponse{
		UserID: u.Id,
		Name:   u.Name,
		Email:  u.Email,
	}, nil
}

func (pq *Postgres) Unlock(ctx context.Context, req *UnlockRequest) (*UnlockResponse, error) {
	u := &User{}
	err := pq.db.Get(u, "SELECT * FROM users WHERE id=$1", req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.State != 4 {
		return nil, errors.New("user is not locked")
	}

	_, err = pq.db.ExecContext(ctx, "UPDATE  users SET  state=2  WHERE id=$1", u.Id)
	if err != nil {
		return nil, err
	}

	return &UnlockResponse{
		UserID: u.Id,
		Name:   u.Name,
		Email:  u.Email,
	}, nil
}

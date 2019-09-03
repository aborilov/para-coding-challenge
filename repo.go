package main

import (
	"database/sql"

	"github.com/doug-martin/goqu/v8"
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"
)

const (
	table = "jwt"
)

type repository struct {
	db *goqu.Database
}

type Repository interface {
	Get(email string) (*User, error)
	Add(u *User) (*User, error)
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: goqu.New("postgres", db),
	}
}

func (r *repository) Get(email string) (*User, error) {
	user := &User{}
	found, err := r.db.From(table).Where(goqu.I("email").Eq(email)).ScanStruct(user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNotFound{Email: email}
	}
	return user, nil
}

func (r *repository) Add(u *User) (*User, error) {
	_, err := r.db.Insert(table).Rows(u).Executor().Exec()
	if err != nil {
		return nil, err
	}
	return r.Get(u.Email)
}

package main

import (
	"database/sql"
	"errors"
)

type user struct {
	ID         int64  `json:"id"`
	Email      string `json:"string"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	CreatedAt  string `json:"createdAt"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}

func (u *user) getUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *user) updateUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *user) deleteUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *user) createUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *user) getUsers(db *sql.DB) error {
	return errors.New("Not implemented")
}

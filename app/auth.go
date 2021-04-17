package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type userLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userRegister struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *userRegister) register(db *sql.DB) error {
	err := u.setPassword()

	if err != nil {
		return err
	}

	err = db.QueryRow(
		"INSERT INTO users(email, password, firstName, lastName) VALUES($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.Password,
		u.FirstName,
		u.LastName,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *userLogin) getUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email, &u.Password)
}

func (u *userRegister) checkUserExists(db *sql.DB) (bool, error) {
	err := db.QueryRow(
		"SELECT email FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email)

	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (u *userRegister) setPassword() error {
	hashingCost := 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashingCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
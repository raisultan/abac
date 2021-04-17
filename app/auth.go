package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type fieldValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type validationError struct {
	Details []fieldValidationError `json:"details"`
}

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginJWTResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type userRegisterRequest struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`

	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type userRegisterResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *userRegisterRequest) register(db *sql.DB) error {
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

func (u *userLoginRequest) getUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email, &u.Password)
}

func (u *userRegisterRequest) checkUserExists(db *sql.DB) (bool, error) {
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

func (u *userRegisterRequest) setPassword() error {
	hashingCost := 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashingCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

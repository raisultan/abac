package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID int `json:"id"`

	Email    string `json:"email"`
	Password string `json:"password"`

	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	CreatedAt  string `json:"createdAt"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}

func (u *user) getUser(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password, firstName, lastName FROM users WHERE id=$1",
		u.ID,
	).Scan(&u.Email, &u.Password, &u.FirstName, &u.LastName)
}

func (u *user) getUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password, firstName, lastName FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email, &u.Password, &u.FirstName, &u.LastName)
}

func (u *user) updateUser(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE users SET email=$1, firstName=$2, lastName=$3 WHERE id=$4",
		u.Email,
		u.FirstName,
		u.LastName,
		u.ID,
	)

	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}

func (u *user) setPassword() error {
	hashingCost := 8
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashingCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *user) createUser(db *sql.DB) error {
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

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	rows, err := db.Query(
		"SELECT id, email, password, firstName, lastName FROM users LIMIT $1 OFFSET $2",
		count,
		start,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

package main

import (
	"database/sql"
)

type user struct {
	ID int `json:"id"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	CreatedAt  string `json:"createdAt"`
	IsAdmin    bool   `json:"isAdmin"`
	IsApproved bool   `json:"isApproved"`
}

type userUpdateRequest struct {
	ID int `json:"-"`

	Email      string `json:"-"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	CreatedAt  string `json:"-"`
	IsAdmin    bool   `json:"-"`
	IsApproved bool   `json:"-"`
}

type userUpdateResponse struct {
	ID int `json:"id"`

	Email      string `json:"email" validate:"required,email"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
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

func (u *userUpdateRequest) updateUser(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE users SET firstName=$1, lastName=$2 WHERE id=$3",
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

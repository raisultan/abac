package main

import "database/sql"

type userLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *userLogin) getUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT email, password FROM users WHERE email=$1",
		u.Email,
	).Scan(&u.Email, &u.Password)
}

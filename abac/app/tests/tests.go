package main_test

import (
	"log"
	"os"
	"testing"

	"abac/main"
)

var a main.App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS users
	(
		id SERIAL,
		email TEXT NOT NULL,
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		createdAt TEXT NOT NULL,
		isAdmin BOOLEAN NOT NULL DEFAULT FALSE,
		isApproved BOOLEAN NOT NULL DEFAULT FALSE,

		CONSTRAINT users_pkey PRIMARY KEY (id)
	)
`

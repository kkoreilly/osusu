package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// ConnectToDB connects to the database
func ConnectToDB() error {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	return err
}

// CreateUserDB creates a user in the database
func CreateUserDB(user User) error {
	statement := `INSERT INTO users (username, password)
	VALUES ($1, $2)`
	_, err := db.Exec(statement, user.Username, user.Password)
	return err
}

// SignInDB checks whether a specific user can sign into the database
func SignInDB(user User) error {
	statement := `SELECT password FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var password string
	err := row.Scan(&password)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

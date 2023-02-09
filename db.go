package main

import (
	"database/sql"
	"errors"

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

// CreateUserDB creates a user in the database and returns the created user if successful and an error if not
func CreateUserDB(user User) (User, error) {
	statement := `INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id`
	row := db.QueryRow(statement, user.Username, user.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return User{}, err
	}
	user.ID = id
	return user, err
}

// SignInDB checks whether a specific user can sign into the database and returns the user if they can and an error if they can't
func SignInDB(user User) (User, error) {
	statement := `SELECT id, password FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var id int
	var password string
	err := row.Scan(&id, &password)
	if err != nil {
		return User{}, errors.New("no user with the given username exists")
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		return User{}, errors.New("incorrect password")
	}
	user.ID = id
	return user, nil
}

// GetMealsDB gets the meals from the database that are owned by the given user
func GetMealsDB(user User) (Meals, error) {
	statement := `SELECT * FROM meals WHERE owner=$1`
	rows, err := db.Query(statement, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(Meals)
	for rows.Next() {
		var id, cost, effort, healthiness, owner int
		var name string
		err := rows.Scan(&id, &name, &cost, &effort, &healthiness, &owner)
		if err != nil {
			return nil, err
		}
		meal := Meal{
			ID:          id,
			Name:        name,
			Cost:        cost,
			Effort:      effort,
			Healthiness: healthiness,
		}
		res[name] = meal
	}
	return res, nil
}

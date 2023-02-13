package main

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
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
	statement := `SELECT id, password, people FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var id int
	var password string
	var people []string
	err := row.Scan(&id, &password, pq.Array(&people))
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("no user with the given username exists")
		}
		return User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return User{}, errors.New("incorrect password")
		}
		return User{}, err
	}
	user.ID = id
	user.People = people
	return user, nil
}

// GetMealsDB gets the meals from the database that are owned by the given user
func GetMealsDB(owner int) (Meals, error) {
	statement := `SELECT * FROM meals WHERE owner=$1`
	rows, err := db.Query(statement, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := Meals{}
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
			Owner:       owner,
		}
		res = append(res, meal)
	}
	return res, nil
}

// CreateMealDB creates a meal in the database with the given owner and returns the created meal if successful and an error if not
func CreateMealDB(owner int) (Meal, error) {
	statement := `INSERT INTO meals (owner)
	VALUES ($1) RETURNING *`
	row := db.QueryRow(statement, owner)
	var id, cost, effort, healthiness int
	var name string
	err := row.Scan(&id, &name, &cost, &effort, &healthiness, &owner)
	if err != nil {
		return Meal{}, err
	}
	meal := Meal{
		ID:          id,
		Name:        name,
		Cost:        cost,
		Effort:      effort,
		Healthiness: healthiness,
		Owner:       owner,
	}
	return meal, nil
}

// UpdateMealDB updates a meal in the database
func UpdateMealDB(meal Meal) error {
	statement := `UPDATE meals
	SET name = $1, cost = $2, effort = $3, healthiness = $4
	WHERE id = $5`
	_, err := db.Exec(statement, meal.Name, meal.Cost, meal.Effort, meal.Healthiness, meal.ID)
	return err
}

// DeleteMealDB deletes the meal with the given id from the database
func DeleteMealDB(id int) error {
	statement := `DELETE FROM meals
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

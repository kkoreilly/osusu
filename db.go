package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

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
	err := row.Scan(&user.ID)
	if err != nil {
		return User{}, err
	}
	return user, err
}

// SignInDB checks whether a specific user can sign into the database and returns the user if they can and an error if they can't
func SignInDB(user User) (User, error) {
	statement := `SELECT id, password FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var password string
	err := row.Scan(&user.ID, &password)
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
	return user, nil
}

// CreateSessionDB creates a new session in the database with the given session id and user id
func CreateSessionDB(id string, userID int) error {
	statement := `INSERT INTO sessions (id, user_id, expires)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(statement, id, userID, time.Now().UTC().Add(RememberMeSessionLength))
	return err
}

// GetSessionDB gets the user id and expiration date of the given session if it exists. Otherwise, it returns an error
func GetSessionDB(id string) (userID int, expires time.Time, err error) {
	statement := `SELECT user_id, expires FROM sessions WHERE id=$1`
	row := db.QueryRow(statement, id)
	err = row.Scan(&userID, &expires)
	return
}

// DeleteSessionDB deletes the given session from the database
func DeleteSessionDB(id string) error {
	statement := `DELETE FROM sessions
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// GetMealsDB gets the meals from the database that are associated with the given user id
func GetMealsDB(userID int) (Meals, error) {
	statement := `SELECT id, name, description, cost, effort, healthiness, taste, type, source, last_done FROM meals WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res Meals
	for rows.Next() {
		var meal Meal
		var tasteJSON string
		err := rows.Scan(&meal.ID, &meal.Name, &meal.Description, &meal.Cost, &meal.Effort, &meal.Healthiness, &tasteJSON, &meal.Type, &meal.Source, &meal.LastDone)
		if err != nil {
			return nil, err
		}
		var tasteMap map[int]int
		err = json.Unmarshal([]byte(tasteJSON), &tasteMap)
		if err != nil {
			return nil, err
		}
		meal.Taste = tasteMap
		res = append(res, meal)
	}
	return res, nil
}

// CreateMealDB creates a meal in the database with the given userID and returns the created meal if successful and an error if not
func CreateMealDB(userID int) (Meal, error) {
	statement := `INSERT INTO meals (user_id, last_done)
	VALUES ($1, $2) RETURNING id, cost, effort, healthiness, taste, type, source`
	row := db.QueryRow(statement, userID, time.Now())
	var meal Meal
	var tasteJSON string
	err := row.Scan(&meal.ID, &meal.Cost, &meal.Effort, &meal.Healthiness, &tasteJSON, &meal.Type, &meal.Source)
	if err != nil {
		return Meal{}, err
	}
	var tasteMap map[int]int
	err = json.Unmarshal([]byte(tasteJSON), &tasteMap)
	if err != nil {
		return Meal{}, err
	}
	meal.Taste = tasteMap
	return meal, nil
}

// UpdateMealDB updates a meal in the database
func UpdateMealDB(meal Meal) error {
	statement := `UPDATE meals
	SET name = $1, description = $2, cost = $3, effort = $4, healthiness = $5, taste = $6, type = $7, source = $8, last_done = $9
	WHERE id = $10`
	tasteJSON, err := json.Marshal(meal.Taste)
	if err != nil {
		return err
	}
	_, err = db.Exec(statement, meal.Name, meal.Description, meal.Cost, meal.Effort, meal.Healthiness, string(tasteJSON), meal.Type, meal.Source, meal.LastDone, meal.ID)
	return err
}

// DeleteMealDB deletes the meal with the given id from the database
func DeleteMealDB(id int) error {
	statement := `DELETE FROM meals
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// GetPeopleDB gets all of the people from the database that are associated with the given user id
func GetPeopleDB(userID int) (People, error) {
	statement := `SELECT id, name FROM people WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res People
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name)
		if err != nil {
			return nil, err
		}
		person.UserID = userID
		res = append(res, person)
	}
	return res, nil
}

// CreatePersonDB creates a person in the database with the given userID and returns the created person if successful and an error if not
func CreatePersonDB(userID int) (Person, error) {
	statement := `INSERT INTO people (user_id)
	VALUES ($1) RETURNING id, name`
	row := db.QueryRow(statement, userID)
	var person Person
	err := row.Scan(&person.ID, &person.Name)
	if err != nil {
		return Person{}, err
	}
	person.UserID = userID
	return person, nil
}

// UpdatePersonDB updates a person in the database
func UpdatePersonDB(person Person) error {
	statement := `UPDATE people
	SET name = $1
	WHERE id = $2`
	_, err := db.Exec(statement, person.Name, person.ID)
	return err
}

// DeletePersonDB deletes the person with the given id from the database
func DeletePersonDB(id int) error {
	statement := `DELETE FROM people
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

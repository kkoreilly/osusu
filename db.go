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
	return user, nil
}

// CreateSessionDB creates a new session in the database with the given session id and user id
func CreateSessionDB(id string, userID int) error {
	statement := `INSERT INTO sessions (id, user_id, expires)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(statement, id, userID, time.Now().UTC().Add(30*24*time.Hour))
	return err
}

// GetSessionDB gets the user id and expiration date of the given session if it exists. Otherwise, it returns an error
func GetSessionDB(id string) (userID int, expires time.Time, err error) {
	statement := `SELECT user_id, expires FROM sessions WHERE id=$1`
	row := db.QueryRow(statement, id)
	err = row.Scan(&userID, &expires)
	return
}

// GetMealsDB gets the meals from the database that are associated with the given user id
func GetMealsDB(userID int) (Meals, error) {
	statement := `SELECT id, name, cost, effort, healthiness, taste, type, source, last_done FROM meals WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := Meals{}
	for rows.Next() {
		var id, cost, effort, healthiness, userID int
		var name, taste, mealType, source string
		var lastDone time.Time
		err := rows.Scan(&id, &name, &cost, &effort, &healthiness, &taste, &mealType, &source, &lastDone)
		if err != nil {
			return nil, err
		}
		var tasteMap map[int]int
		err = json.Unmarshal([]byte(taste), &tasteMap)
		if err != nil {
			return nil, err
		}
		meal := Meal{
			ID:          id,
			Name:        name,
			Cost:        cost,
			Effort:      effort,
			Healthiness: healthiness,
			Taste:       tasteMap,
			Type:        mealType,
			Source:      source,
			LastDone:    lastDone,
			UserID:      userID,
		}
		res = append(res, meal)
	}
	return res, nil
}

// CreateMealDB creates a meal in the database with the given userID and returns the created meal if successful and an error if not
func CreateMealDB(userID int) (Meal, error) {
	statement := `INSERT INTO meals (user_id, last_done)
	VALUES ($1, $2) RETURNING id, name, cost, effort, healthiness, taste, type, source`
	row := db.QueryRow(statement, userID, time.Now())
	var id, cost, effort, healthiness int
	var name, taste, mealType, source string
	err := row.Scan(&id, &name, &cost, &effort, &healthiness, &taste, &mealType, &source)
	if err != nil {
		return Meal{}, err
	}
	var tasteMap map[int]int
	err = json.Unmarshal([]byte(taste), &tasteMap)
	if err != nil {
		return Meal{}, err
	}
	meal := Meal{
		ID:          id,
		Name:        name,
		Cost:        cost,
		Effort:      effort,
		Healthiness: healthiness,
		Taste:       tasteMap,
		Type:        mealType,
		Source:      source,
		LastDone:    time.Now(),
		UserID:      userID,
	}
	return meal, nil
}

// UpdateMealDB updates a meal in the database
func UpdateMealDB(meal Meal) error {
	statement := `UPDATE meals
	SET name = $1, cost = $2, effort = $3, healthiness = $4, taste = $5, type = $6, source = $7, last_done = $8
	WHERE id = $9`
	tasteJSON, err := json.Marshal(meal.Taste)
	if err != nil {
		return err
	}
	_, err = db.Exec(statement, meal.Name, meal.Cost, meal.Effort, meal.Healthiness, string(tasteJSON), meal.Type, meal.Source, meal.LastDone, meal.ID)
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
	statement := `SELECT * FROM people WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := People{}
	for rows.Next() {
		var id, userID int
		var name string
		err := rows.Scan(&id, &name, &userID)
		if err != nil {
			return nil, err
		}
		person := Person{
			ID:     id,
			Name:   name,
			UserID: userID,
		}
		res = append(res, person)
	}
	return res, nil
}

// CreatePersonDB creates a person in the database with the given userID and returns the created person if successful and an error if not
func CreatePersonDB(userID int) (Person, error) {
	statement := `INSERT INTO people (user_id)
	VALUES ($1) RETURNING *`
	row := db.QueryRow(statement, userID)
	var id int
	var name string
	err := row.Scan(&id, &name, &userID)
	if err != nil {
		return Person{}, err
	}
	person := Person{
		ID:     id,
		Name:   name,
		UserID: userID,
	}
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

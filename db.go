package main

import (
	"database/sql"
	"encoding/json"
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

// GetMealsDB gets the meals from the database that are associated with the given user id
func GetMealsDB(userID int) (Meals, error) {
	statement := `SELECT * FROM meals WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := Meals{}
	for rows.Next() {
		var id, cost, effort, healthiness, userID int
		var name, taste string
		err := rows.Scan(&id, &name, &cost, &effort, &healthiness, &taste, &userID)
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
			UserID:      userID,
		}
		res = append(res, meal)
	}
	return res, nil
}

// CreateMealDB creates a meal in the database with the given userID and returns the created meal if successful and an error if not
func CreateMealDB(userID int) (Meal, error) {
	statement := `INSERT INTO meals (user_id)
	VALUES ($1) RETURNING *`
	row := db.QueryRow(statement, userID)
	var id, cost, effort, healthiness int
	var name string
	err := row.Scan(&id, &name, &cost, &effort, &healthiness, &userID)
	if err != nil {
		return Meal{}, err
	}
	meal := Meal{
		ID:          id,
		Name:        name,
		Cost:        cost,
		Effort:      effort,
		Healthiness: healthiness,
		UserID:      userID,
	}
	return meal, nil
}

// UpdateMealDB updates a meal in the database
func UpdateMealDB(meal Meal) error {
	statement := `UPDATE meals
	SET name = $1, cost = $2, effort = $3, healthiness = $4, taste = $5
	WHERE id = $6`
	tasteJSON, err := json.Marshal(meal.Taste)
	if err != nil {
		return err
	}
	_, err = db.Exec(statement, meal.Name, meal.Cost, meal.Effort, meal.Healthiness, string(tasteJSON), meal.ID)
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

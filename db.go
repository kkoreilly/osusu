package main

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// ConnectToDB connects to the database
func ConnectToDB() error {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return err
}

// CreateUserDB creates a user in the database and returns the created user if successful and an error if not
func CreateUserDB(user User) (User, error) {
	statement := `INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id, cuisines`
	row := db.QueryRow(statement, user.Username, user.Password)
	err := row.Scan(&user.ID, pq.Array(&user.Cuisines))
	if err != nil {
		return User{}, err
	}
	return user, err
}

// SignInDB checks whether a specific user can sign into the database and returns the user if they can and an error if they can't
func SignInDB(user User) (User, error) {
	statement := `SELECT id, password, cuisines FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var password string
	err := row.Scan(&user.ID, &password, pq.Array(&user.Cuisines))
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

// GetUserCuisinesDB gets the cuisines value of the user with the given id from the database
func GetUserCuisinesDB(userID int64) ([]string, error) {
	statement := `SELECT cuisines FROM users WHERE id = $1`
	row := db.QueryRow(statement, userID)
	var cuisines []string
	err := row.Scan(pq.Array(&cuisines))
	if err != nil {
		return nil, err
	}
	return cuisines, nil
}

// UpdateUserCuisinesDB updates the cuisines value of the given user in the database to its cuisines value
func UpdateUserCuisinesDB(user User) error {
	statement := `UPDATE users
	SET cuisines = $1
	WHERE id = $2`
	_, err := db.Exec(statement, pq.Array(user.Cuisines), user.ID)
	return err
}

// UpdateUsernameDB updates the username of the given user in the database to its username value
func UpdateUsernameDB(user User) error {
	statement := `UPDATE users
	SET username = $1
	WHERE id = $2`
	_, err := db.Exec(statement, user.Username, user.ID)
	return err
}

// UpdatePasswordDB updates the password of the given user in the database to its password value
func UpdatePasswordDB(user User) error {
	statement := `UPDATE users
	SET password = $1
	WHERE id = $2`
	_, err := db.Exec(statement, user.Password, user.ID)
	return err
}

// CreateSessionDB creates a new session in the database with the given session id and user id
func CreateSessionDB(id string, userID int64) error {
	statement := `INSERT INTO sessions (id, user_id, expires)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(statement, id, userID, time.Now().UTC().Add(RememberMeSessionLength))
	return err
}

// GetSessionDB gets the user id and expiration date of the given session if it exists. Otherwise, it returns an error
func GetSessionDB(id string) (userID int64, expires time.Time, err error) {
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

// GetGroupsDB returns the groups that the user with the given user id is part of
func GetGroupsDB(userID int64) (Groups, error) {
	statement := `SELECT id, owner, code, name, members FROM groups WHERE $1 = ANY (members)`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	res := Groups{}
	for rows.Next() {
		group := Group{}
		err := rows.Scan(&group.ID, &group.Owner, &group.Code, &group.Name, &group.Members)
		if err != nil {
			return nil, err
		}
		res = append(res, group)
	}
	return res, nil
}

// CreateGroupDB creates the given group in the database and returns the created group if successful and an error if not
func CreateGroupDB(group Group) (Group, error) {
	statement := `INSERT INTO groups (owner, code, name, members)
	VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(statement, group.Owner, group.Code, group.Name, pq.Array(group.Members))
	err := row.Scan(&group.ID)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

// JoinGroupDB has the user with the given user id join the group with the given group code and returns the joined group if successful and an error if not
func JoinGroupDB(groupJoin GroupJoin) (Group, error) {
	statement := `UPDATE groups
	SET members = ARRAY_APPEND(members, $1)
	WHERE code = $2
	RETURNING id, owner, name, members`
	row := db.QueryRow(statement, groupJoin.UserID, groupJoin.GroupCode)
	res := Group{}
	err := row.Scan(&res.ID, &res.Owner, &res.Name, &res.Members)
	if err != nil {
		return Group{}, err
	}
	return res, nil
}

// UpdateGroupDB updates the name and members of the group with the given id to the values of the given group
func UpdateGroupDB(group Group) error {
	statement := `UPDATE groups
	SET name = $1, members = $2
	WHERE id = $3`
	_, err := db.Exec(statement, group.Name, group.Members)
	return err
}

// GetMealsDB gets the meals from the database that are associated with the given user id
func GetMealsDB(userID int64) (Meals, error) {
	statement := `SELECT id, name, description, cuisine FROM meals WHERE user_id=$1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res Meals
	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.ID, &meal.Name, &meal.Description, pq.Array(&meal.Cuisine))
		if err != nil {
			return nil, err
		}
		res = append(res, meal)
	}
	return res, nil
}

// CreateMealDB creates a meal in the database with the given information and returns the created meal if successful and an error if not
func CreateMealDB(meal Meal) (Meal, error) {
	statement := `INSERT INTO meals (user_id, name, description, cuisine)
	VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(statement, meal.UserID, meal.Name, meal.Description, pq.Array(meal.Cuisine))
	err := row.Scan(&meal.ID)
	if err != nil {
		return Meal{}, err
	}
	return meal, nil
}

// UpdateMealDB updates a meal in the database
func UpdateMealDB(meal Meal) error {
	statement := `UPDATE meals
	SET name = $1, description = $2, cuisine = $3
	WHERE id = $4`
	_, err := db.Exec(statement, meal.Name, meal.Description, pq.Array(meal.Cuisine), meal.ID)
	return err
}

// DeleteMealDB deletes the meal with the given id from the database
func DeleteMealDB(id int64) error {
	statement := `DELETE FROM meals
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// DeleteMealEntriesDB deletes the entries associated with the given meal id from the database
func DeleteMealEntriesDB(mealID int64) error {
	statement := `DELETE FROM entries
	WHERE meal_id = $1`
	_, err := db.Exec(statement, mealID)
	return err
}

// GetPeopleDB gets all of the people from the database that are associated with the given user id
func GetPeopleDB(userID int64) (People, error) {
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
func CreatePersonDB(userID int64) (Person, error) {
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
func DeletePersonDB(id int64) error {
	statement := `DELETE FROM people
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// GetEntriesDB gets the entries from the database that have the given user id
func GetEntriesDB(userID int64) (Entries, error) {
	statement := `SELECT id, meal_id, entry_date, type, source, cost, effort, healthiness, taste FROM entries
	WHERE user_id = $1`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	var res Entries
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.MealID, &entry.Date, &entry.Type, &entry.Source, &entry.Cost, &entry.Effort, &entry.Healthiness, &entry.Taste)
		if err != nil {
			return nil, err
		}
		entry.UserID = userID
		res = append(res, entry)
	}
	return res, nil
}

// GetEntriesForMealDB gets the entries from the database that have the given meal id
func GetEntriesForMealDB(mealID int64) (Entries, error) {
	statement := `SELECT id, user_id, entry_date, type, source, cost, effort, healthiness, taste FROM entries
	WHERE meal_id = $1`
	rows, err := db.Query(statement, mealID)
	if err != nil {
		return nil, err
	}
	var res Entries
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Type, &entry.Source, &entry.Cost, &entry.Effort, &entry.Healthiness, &entry.Taste)
		if err != nil {
			return nil, err
		}
		entry.MealID = mealID
		res = append(res, entry)
	}
	return res, nil
}

// CreateEntryDB creates and returns a new entry in the database with the given entry's user and meal id values
func CreateEntryDB(entry Entry) (Entry, error) {
	statement := `INSERT INTO entries (user_id, meal_id, entry_date, type, source, cost, effort, healthiness, taste)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := db.QueryRow(statement, entry.UserID, entry.MealID, entry.Date, entry.Type, entry.Source, entry.Cost, entry.Effort, entry.Healthiness, entry.Taste)
	err := row.Scan(&entry.ID)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

// UpdateEntryDB updates an entry in the database to have the values of the given entry
func UpdateEntryDB(entry Entry) error {
	statement := `UPDATE entries
	SET entry_date = $1, type = $2, source = $3, cost = $4, effort = $5, healthiness = $6, taste = $7
	WHERE id = $8`
	_, err := db.Exec(statement, entry.Date, entry.Type, entry.Source, entry.Cost, entry.Effort, entry.Healthiness, entry.Taste, entry.ID)
	return err
}

// DeleteEntryDB deletes the entry with the given id from the database
func DeleteEntryDB(id int64) error {
	statement := `DELETE FROM entries
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

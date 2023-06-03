package db

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/kkoreilly/osusu/osusu"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// db is the database instance used in the app
var db *sql.DB

// Start connects to, initializes, and cleans up the database. It should be called on server start.
func Start() error {
	err := Connect()
	if err != nil {
		return err
	}
	err = Init()
	if err != nil {
		return err
	}
	err = Cleanup()
	if err != nil {
		return err
	}
	return nil
}

// Connect connects to the database
func Connect() error {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return err
}

// CreateUser creates a user in the database and returns the created user if successful and an error if not
func CreateUser(user osusu.User) (osusu.User, error) {
	statement := `INSERT INTO users (username, password, name)
	VALUES ($1, $2, $3) RETURNING id`
	row := db.QueryRow(statement, user.Username, user.Password, user.Name)
	err := row.Scan(&user.ID)
	if err != nil {
		return osusu.User{}, err
	}
	return user, err
}

// SignIn checks whether a specific user can sign into the database and returns the user if they can and an error if they can't
func SignIn(user osusu.User) (osusu.User, error) {
	statement := `SELECT id, password, name FROM users WHERE username=$1`
	row := db.QueryRow(statement, user.Username)
	var password string
	err := row.Scan(&user.ID, &password, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return osusu.User{}, errors.New("no user with the given username exists")
		}
		return osusu.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return osusu.User{}, errors.New("incorrect password")
		}
		return osusu.User{}, err
	}
	return user, nil
}

// GetGroupCuisines gets the cuisines value of the group with the given id from the database
func GetGroupCuisines(groupID int64) ([]string, error) {
	statement := `SELECT cuisines FROM groups WHERE id = $1`
	row := db.QueryRow(statement, groupID)
	var cuisines []string
	err := row.Scan(pq.Array(&cuisines))
	if err != nil {
		return nil, err
	}
	return cuisines, nil
}

// UpdateGroupCuisines updates the cuisines value of the given group in the database to its cuisines value
func UpdateGroupCuisines(group osusu.Group) error {
	statement := `UPDATE groups
	SET cuisines = $1
	WHERE id = $2`
	_, err := db.Exec(statement, pq.Array(group.Cuisines), group.ID)
	return err
}

// UpdateUserInfo updates the username and name of the given user in the database
func UpdateUserInfo(user osusu.User) error {
	statement := `UPDATE users
	SET username = $1, name = $2
	WHERE id = $3`
	_, err := db.Exec(statement, user.Username, user.Name, user.ID)
	return err
}

// UpdatePassword updates the password of the given user in the database to its password value
func UpdatePassword(user osusu.User) error {
	statement := `UPDATE users
	SET password = $1
	WHERE id = $2`
	_, err := db.Exec(statement, user.Password, user.ID)
	return err
}

// GetUsers gets all of the users with user ids contained within the given array of user ids
func GetUsers(userIDs []int64) (osusu.Users, error) {
	statement := `SELECT id, username, name FROM users WHERE id = ANY($1)`
	rows, err := db.Query(statement, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	res := osusu.Users{}
	for rows.Next() {
		var user osusu.User
		err := rows.Scan(&user.ID, &user.Username, &user.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}

// CreateSession creates a new session in the database with the given session id and user id
func CreateSession(id string, userID int64) error {
	statement := `INSERT INTO sessions (id, user_id, expires)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(statement, id, userID, time.Now().UTC().Add(osusu.RememberMeSessionLength))
	return err
}

// GetSession gets the user id and expiration date of the given session if it exists. Otherwise, it returns an error
func GetSession(id string) (userID int64, expires time.Time, err error) {
	statement := `SELECT user_id, expires FROM sessions WHERE id=$1`
	row := db.QueryRow(statement, id)
	err = row.Scan(&userID, &expires)
	return
}

// DeleteSession deletes the given session from the database
func DeleteSession(id string) error {
	statement := `DELETE FROM sessions
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// GetGroups returns the groups that the user with the given user id is part of
func GetGroups(userID int64) (osusu.Groups, error) {
	statement := `SELECT id, owner, code, name, members FROM groups WHERE $1 = ANY (members)`
	rows, err := db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	res := osusu.Groups{}
	for rows.Next() {
		group := osusu.Group{}
		err := rows.Scan(&group.ID, &group.Owner, &group.Code, &group.Name, pq.Array(&group.Members))
		if err != nil {
			return nil, err
		}
		res = append(res, group)
	}
	return res, nil
}

// CreateGroup creates the given group in the database and returns the created group if successful and an error if not
func CreateGroup(group osusu.Group) (osusu.Group, error) {
	statement := `INSERT INTO groups (owner, code, name, members)
	VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(statement, group.Owner, group.Code, group.Name, pq.Array(group.Members))
	err := row.Scan(&group.ID)
	if err != nil {
		return osusu.Group{}, err
	}
	return group, nil
}

// JoinGroup has the user with the given user id join the group with the given group code and returns the joined group if successful and an error if not
func JoinGroup(groupJoin osusu.GroupJoin) (osusu.Group, error) {
	statement := `UPDATE groups
	SET members = ARRAY_APPEND(members, $1)
	WHERE code = $2
	RETURNING id, owner, name, members`
	row := db.QueryRow(statement, groupJoin.UserID, groupJoin.GroupCode)
	res := osusu.Group{}
	err := row.Scan(&res.ID, &res.Owner, &res.Name, pq.Array(&res.Members))
	if err != nil {
		return osusu.Group{}, err
	}
	res.Code = groupJoin.GroupCode
	return res, nil
}

// UpdateGroup updates the name and members of the group with the given id to the values of the given group
func UpdateGroup(group osusu.Group) error {
	statement := `UPDATE groups
	SET name = $1, members = $2
	WHERE id = $3`
	_, err := db.Exec(statement, group.Name, pq.Array(group.Members), group.ID)
	return err
}

// GetMeals gets the meals from the database that are associated with the given group id
func GetMeals(groupID int64) (osusu.Meals, error) {
	statement := `SELECT id, name, description, source, image, category, cuisine FROM meals WHERE group_id=$1`
	rows, err := db.Query(statement, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res osusu.Meals
	for rows.Next() {
		var meal osusu.Meal
		err := rows.Scan(&meal.ID, &meal.Name, &meal.Description, &meal.Source, &meal.Image, pq.Array(&meal.Category), pq.Array(&meal.Cuisine))
		if err != nil {
			return nil, err
		}
		res = append(res, meal)
	}
	return res, nil
}

// CreateMeal creates a meal in the database with the given information and returns the created meal if successful and an error if not
func CreateMeal(meal osusu.Meal) (osusu.Meal, error) {
	statement := `INSERT INTO meals (group_id, name, description, source, image, category, cuisine)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	row := db.QueryRow(statement, meal.GroupID, meal.Name, meal.Description, meal.Source, meal.Image, pq.Array(meal.Category), pq.Array(meal.Cuisine))
	err := row.Scan(&meal.ID)
	if err != nil {
		return osusu.Meal{}, err
	}
	return meal, nil
}

// UpdateMeal updates a meal in the database
func UpdateMeal(meal osusu.Meal) error {
	statement := `UPDATE meals
	SET name = $1, description = $2, source = $3, image = $4, category = $5, cuisine = $6
	WHERE id = $7`
	_, err := db.Exec(statement, meal.Name, meal.Description, meal.Source, meal.Image, pq.Array(meal.Category), pq.Array(meal.Cuisine), meal.ID)
	return err
}

// DeleteMeal deletes the meal with the given id from the database
func DeleteMeal(id int64) error {
	statement := `DELETE FROM meals
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

// DeleteMealEntries deletes the entries associated with the given meal id from the database
func DeleteMealEntries(mealID int64) error {
	statement := `DELETE FROM entries
	WHERE meal_id = $1`
	_, err := db.Exec(statement, mealID)
	return err
}

// GetEntries gets the entries from the database that have the given group id
func GetEntries(groupID int64) (osusu.Entries, error) {
	statement := `SELECT id, meal_id, entry_date, type, source, cost, effort, healthiness, taste FROM entries
	WHERE group_id = $1`
	rows, err := db.Query(statement, groupID)
	if err != nil {
		return nil, err
	}
	var res osusu.Entries
	for rows.Next() {
		var entry osusu.Entry
		err := rows.Scan(&entry.ID, &entry.MealID, &entry.Date, &entry.Type, &entry.Source, &entry.Cost, &entry.Effort, &entry.Healthiness, &entry.Taste)
		if err != nil {
			return nil, err
		}
		entry.GroupID = groupID
		res = append(res, entry)
	}
	return res, nil
}

// GetEntriesForMeal gets the entries from the database that have the given meal id
func GetEntriesForMeal(mealID int64) (osusu.Entries, error) {
	statement := `SELECT id, group_id, entry_date, type, source, cost, effort, healthiness, taste FROM entries
	WHERE meal_id = $1`
	rows, err := db.Query(statement, mealID)
	if err != nil {
		return nil, err
	}
	var res osusu.Entries
	for rows.Next() {
		var entry osusu.Entry
		err := rows.Scan(&entry.ID, &entry.GroupID, &entry.Date, &entry.Type, &entry.Source, &entry.Cost, &entry.Effort, &entry.Healthiness, &entry.Taste)
		if err != nil {
			return nil, err
		}
		entry.MealID = mealID
		res = append(res, entry)
	}
	return res, nil
}

// CreateEntry creates and returns a new entry in the database with the given entry's user and meal id values
func CreateEntry(entry osusu.Entry) (osusu.Entry, error) {
	statement := `INSERT INTO entries (group_id, meal_id, entry_date, type, source, cost, effort, healthiness, taste)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := db.QueryRow(statement, entry.GroupID, entry.MealID, entry.Date, entry.Type, entry.Source, entry.Cost, entry.Effort, entry.Healthiness, entry.Taste)
	err := row.Scan(&entry.ID)
	if err != nil {
		return osusu.Entry{}, err
	}
	return entry, nil
}

// UpdateEntry updates an entry in the database to have the values of the given entry
func UpdateEntry(entry osusu.Entry) error {
	statement := `UPDATE entries
	SET entry_date = $1, type = $2, source = $3, cost = $4, effort = $5, healthiness = $6, taste = $7
	WHERE id = $8`
	_, err := db.Exec(statement, entry.Date, entry.Type, entry.Source, entry.Cost, entry.Effort, entry.Healthiness, entry.Taste, entry.ID)
	return err
}

// DeleteEntry deletes the entry with the given id from the database
func DeleteEntry(id int64) error {
	statement := `DELETE FROM entries
	WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

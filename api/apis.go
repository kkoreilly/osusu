package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/kkoreilly/osusu/db"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/kkoreilly/osusu/server"
	"golang.org/x/crypto/bcrypt"
)

// SignUp attempts to create a user with the given information.
// It returns the created user if successful and an error if not.
var SignUp = New(http.MethodPost, "/api/signUp", func(user osusu.User) (osusu.User, error) {
	if user.Username == "" || user.Password == "" || user.Name == "" {
		return osusu.User{}, errors.New("username, password, and name must not be empty")
	}
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return osusu.User{}, err
	}
	user.Password = string(hash)

	user, err = db.CreateUser(user)
	if err != nil {
		return osusu.User{}, err
	}
	// make session if remember me
	if user.RememberMe {
		session, err := server.GenerateSessionID()
		if err != nil {
			return osusu.User{}, err
		}
		user.Session = session
		sessionHash := sha256.Sum256([]byte(session))
		err = db.CreateSession(hex.EncodeToString(sessionHash[:]), user.ID)
		if err != nil {
			return osusu.User{}, err
		}
	}
	// return with blank password
	user.Password = ""
	return user, nil
})

// SignIn attempts to sign the given user into the app.
// It returns the user if successful and an error if not.
var SignIn = New(http.MethodPost, "/api/signIn", func(user osusu.User) (osusu.User, error) {
	user, err := db.SignIn(user)
	if err != nil {
		return osusu.User{}, err
	}
	// make session if remember me
	if user.RememberMe {
		session, err := server.GenerateSessionID()
		if err != nil {
			return osusu.User{}, err
		}
		user.Session = session
		sessionHash := sha256.Sum256([]byte(session))
		err = db.CreateSession(hex.EncodeToString(sessionHash[:]), user.ID)
		if err != nil {
			return osusu.User{}, err
		}
	}
	// return with blank password
	user.Password = ""
	return user, nil
})

// GetGroupCuisines gets the cuisines of the group with the given group id.
var GetGroupCuisines = New(http.MethodGet, "/api/getGroupCuisines", func(groupID int64) ([]string, error) {
	return db.GetGroupCuisines(groupID)
})

// UpdateGroupCuisines updates the cuisines of the given group to be the cuisines value of the provided group.
// It returns a confirmation string if successful and an error if not.
var UpdateGroupCuisines = New(http.MethodPut, "/api/updateGroupCuisines", func(group osusu.Group) (string, error) {
	err := db.UpdateGroupCuisines(group)
	if err != nil {
		return "", err
	}
	return "updated group cuisines", nil
})

// UpdateUserInfo updates the username and name of the given user to the values of the provided user.
// It returns a confirmation string if successful and an error if not.
var UpdateUserInfo = New(http.MethodPut, "/api/updateUserInfo", func(user osusu.User) (string, error) {
	if user.Username == "" || user.Name == "" {
		return "", errors.New("username and name must not be empty")
	}
	err := db.UpdateUserInfo(user)
	if err != nil {
		return "", err
	}
	return "updated user info", nil
})

// UpdatePassword updates the password of the given user to be the password value of the provided user.
// It returns a confirmation string if successful and an error if not.
var UpdatePassword = New(http.MethodPut, "/api/updatePassword", func(user osusu.User) (string, error) {
	if user.Password == "" {
		return "", errors.New("password must not be empty")
	}
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)
	err = db.UpdatePassword(user)
	if err != nil {
		return "", err
	}
	return "updated password", nil
})

// GetUsers returns all of the users with user ids contained within the given array of user ids
var GetUsers = New(http.MethodGet, "/api/getUsers", func(userIDs []int64) (osusu.Users, error) {
	return db.GetUsers(userIDs)
})

// AuthenticateSession checks whether the given user has a valid session id.
// It returns a confirmation string if so and an error if not.
var AuthenticateSession = New(http.MethodPost, "/api/authenticateSession", func(user osusu.User) (string, error) {
	sessionHash := sha256.Sum256([]byte(user.Session))
	sessionHashString := hex.EncodeToString(sessionHash[:])
	userID, expires, err := db.GetSession(sessionHashString)
	if err != nil {
		return "", err
	}
	if expires.Before(time.Now()) {
		err := db.DeleteSession(sessionHashString)
		if err != nil {
			log.Println("error deleting expired Session ID:", err)
		}
		return "", errors.New("session id expired")
	}
	if userID != user.ID {
		return "", errors.New("invalid session id")
	}
	return "authenticated", nil
})

// SignOut attempts to sign the given user out of the app.
// It returns a confirmation string if successful and an error if not.
var SignOut = New(http.MethodDelete, "/api/signOut", func(user osusu.User) (string, error) {
	sessionHash := sha256.Sum256([]byte(user.Session))
	sessionHashString := hex.EncodeToString(sessionHash[:])
	err := db.DeleteSession(sessionHashString)
	if err != nil {
		return "", err
	}
	return "signed out", nil
})

// GetGroups gets the groups that the user with the given user id is part of. It returns the groups if successful and an error if not.
var GetGroups = New(http.MethodGet, "/api/getGroups", func(userID int64) (osusu.Groups, error) {
	return db.GetGroups(userID)
})

// CreateGroup creates a new group in the database with the values of the given group and returns the created group if successful and an error if not.
// The group code is generated by the API on the server.
var CreateGroup = New(http.MethodPost, "/api/createGroup", func(group osusu.Group) (osusu.Group, error) {
	code, err := server.GenerateGroupCode()
	if err != nil {
		return osusu.Group{}, err
	}
	group.Code = code
	group, err = db.CreateGroup(group)
	if err != nil {
		return osusu.Group{}, err
	}
	return group, nil
})

// JoinGroup attempts to have the user with the given user id join the group with the given group code. It gets these values from the given GroupJoin struct.
// It returns the joined group if successful and an error if not
var JoinGroup = New(http.MethodPost, "/api/joinGroup", func(groupJoin osusu.GroupJoin) (osusu.Group, error) {
	return db.JoinGroup(groupJoin)
})

// UpdateGroup updates the name and members of the given group with the given values. It returns a confirmation string if successful and an error if not.
var UpdateGroup = New(http.MethodPut, "/api/updateGroup", func(group osusu.Group) (string, error) {
	err := db.UpdateGroup(group)
	if err != nil {
		return "", err
	}
	return "group updated", nil
})

// GetMeals returns the meals in the database that have the given group id.
var GetMeals = New(http.MethodGet, "/api/getMeals", func(groupID int64) (osusu.Meals, error) {
	return db.GetMeals(groupID)
})

// CreateMeal creates and returns a new meal with the given information.
var CreateMeal = New(http.MethodPost, "/api/createMeal", func(meal osusu.Meal) (osusu.Meal, error) {
	return db.CreateMeal(meal)
})

// UpdateMeal updates the given meal in the database to have the values of the given meal.
// It returns a confirmation string if successful and an error otherwise.
var UpdateMeal = New(http.MethodPut, "/api/updateMeal", func(meal osusu.Meal) (string, error) {
	err := db.UpdateMeal(meal)
	if err != nil {
		return "", err
	}
	return "meal updated", nil
})

// DeleteMeal deletes the given meal from the database.
// It returns a confirmation string if successful and an error if not.
var DeleteMeal = New(http.MethodDelete, "/api/deleteMeal", func(id int64) (string, error) {
	err := db.DeleteMealEntries(id)
	if err != nil {
		return "", err
	}
	err = db.DeleteMeal(id)
	if err != nil {
		return "", err
	}
	return "meal deleted", nil
})

// GetEntries fetches and returns the entries associated with the given group id from the database
var GetEntries = New(http.MethodGet, "/api/getEntries", func(groupID int64) (osusu.Entries, error) {
	return db.GetEntries(groupID)
})

// GetEntriesForMeal fetches and returns the entries associated with the given meal id from the database
var GetEntriesForMeal = New(http.MethodGet, "/api/getEntriesForMeal", func(mealID int64) (osusu.Entries, error) {
	return db.GetEntriesForMeal(mealID)
})

// CreateEntry creates and returns a new entry with the given entry's meal and user id values
var CreateEntry = New(http.MethodPost, "/api/createEntry", func(entry osusu.Entry) (osusu.Entry, error) {
	return db.CreateEntry(entry)
})

// UpdateEntry updates the given entry in the database to have the given information.
// It returns a confirmation string if successful and an error if not.
var UpdateEntry = New(http.MethodPut, "/api/updateEntry", func(entry osusu.Entry) (string, error) {
	err := db.UpdateEntry(entry)
	if err != nil {
		return "", err
	}
	return "entry updated", nil
})

// DeleteEntry deletes the entry with the given id from the database.
// It returns a confirmation string if successful and an error if not.
var DeleteEntry = New(http.MethodDelete, "/api/deleteEntry", func(id int64) (string, error) {
	err := db.DeleteEntry(id)
	if err != nil {
		return "", err
	}
	return "entry deleted", nil
})

// RecommendRecipes returns a list of recommended recipes based on the given word score map
var RecommendRecipes = New(http.MethodGet, "/api/recommendRecipes", func(data osusu.RecommendRecipesData) (osusu.Recipes, error) {
	return server.RecommendRecipes(data), nil
})

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignUpAPI attempts to create a user with the given information.
// It returns the created user if successful and an error if not.
var SignUpAPI = NewAPI(http.MethodPost, "/api/signUp", func(user User) (User, error) {
	if user.Username == "" || user.Password == "" || user.Name == "" {
		return User{}, errors.New("username, password, and name must not be empty")
	}
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return User{}, err
	}
	user.Password = string(hash)

	user, err = CreateUserDB(user)
	if err != nil {
		return User{}, err
	}
	// make session if remember me
	if user.RememberMe {
		session, err := GenerateSessionID()
		if err != nil {
			return User{}, err
		}
		user.Session = session
		sessionHash := sha256.Sum256([]byte(session))
		err = CreateSessionDB(hex.EncodeToString(sessionHash[:]), user.ID)
		if err != nil {
			return User{}, err
		}
	}
	// return with blank password
	user.Password = ""
	return user, nil
})

// SignInAPI attempts to sign the given user into the app.
// It returns the user if successful and an error if not.
var SignInAPI = NewAPI(http.MethodPost, "/api/signIn", func(user User) (User, error) {
	user, err := SignInDB(user)
	if err != nil {
		return User{}, err
	}
	// make session if remember me
	if user.RememberMe {
		session, err := GenerateSessionID()
		if err != nil {
			return User{}, err
		}
		user.Session = session
		sessionHash := sha256.Sum256([]byte(session))
		err = CreateSessionDB(hex.EncodeToString(sessionHash[:]), user.ID)
		if err != nil {
			return User{}, err
		}
	}
	// return with blank password
	user.Password = ""
	return user, nil
})

// GetGroupCuisinesAPI gets the cuisines of the group with the given group id.
var GetGroupCuisinesAPI = NewAPI(http.MethodGet, "/api/getGroupCuisines", func(groupID int64) ([]string, error) {
	return GetGroupCuisinesDB(groupID)
})

// UpdateGroupCuisinesAPI updates the cuisines of the given group to be the cuisines value of the provided group.
// It returns a confirmation string if successful and an error if not.
var UpdateGroupCuisinesAPI = NewAPI(http.MethodPut, "/api/updateGroupCuisines", func(group Group) (string, error) {
	err := UpdateGroupCuisinesDB(group)
	if err != nil {
		return "", err
	}
	return "updated group cuisines", nil
})

// UpdateUserInfoAPI updates the username and name of the given user to the values of the provided user.
// It returns a confirmation string if successful and an error if not.
var UpdateUserInfoAPI = NewAPI(http.MethodPut, "/api/updateUserInfo", func(user User) (string, error) {
	if user.Username == "" || user.Name == "" {
		return "", errors.New("username and name must not be empty")
	}
	err := UpdateUserInfoDB(user)
	if err != nil {
		return "", err
	}
	return "updated user info", nil
})

// UpdatePasswordAPI updates the password of the given user to be the password value of the provided user.
// It returns a confirmation string if successful and an error if not.
var UpdatePasswordAPI = NewAPI(http.MethodPut, "/api/updatePassword", func(user User) (string, error) {
	if user.Password == "" {
		return "", errors.New("password must not be empty")
	}
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return "", err
	}
	user.Password = string(hash)
	err = UpdatePasswordDB(user)
	if err != nil {
		return "", err
	}
	return "updated password", nil
})

// GetUsersAPI returns all of the users with user ids contained within the given array of user ids
var GetUsersAPI = NewAPI(http.MethodGet, "/api/getUsers", func(userIDs []int64) (Users, error) {
	return GetUsersDB(userIDs)
})

// AuthenticateSessionAPI checks whether the given user has a valid session id.
// It returns a confirmation string if so and an error if not.
var AuthenticateSessionAPI = NewAPI(http.MethodPost, "/api/authenticateSession", func(user User) (string, error) {
	sessionHash := sha256.Sum256([]byte(user.Session))
	sessionHashString := hex.EncodeToString(sessionHash[:])
	userID, expires, err := GetSessionDB(sessionHashString)
	if err != nil {
		return "", err
	}
	if expires.Before(time.Now()) {
		err := DeleteSessionDB(sessionHashString)
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

// SignOutAPI attempts to sign the given user out of the app.
// It returns a confirmation string if successful and an error if not.
var SignOutAPI = NewAPI(http.MethodDelete, "/api/signOut", func(user User) (string, error) {
	sessionHash := sha256.Sum256([]byte(user.Session))
	sessionHashString := hex.EncodeToString(sessionHash[:])
	err := DeleteSessionDB(sessionHashString)
	if err != nil {
		return "", err
	}
	return "signed out", nil
})

// GetGroupsAPI gets the groups that the user with the given user id is part of. It returns the groups if successful and an error if not.
var GetGroupsAPI = NewAPI(http.MethodGet, "/api/getGroups", func(userID int64) (Groups, error) {
	return GetGroupsDB(userID)
})

// CreateGroupAPI creates a new group in the database with the values of the given group and returns the created group if successful and an error if not.
// The group code is generated by the API on the server.
var CreateGroupAPI = NewAPI(http.MethodPost, "/api/createGroup", func(group Group) (Group, error) {
	code, err := GenerateGroupCode()
	if err != nil {
		return Group{}, err
	}
	group.Code = code
	group, err = CreateGroupDB(group)
	if err != nil {
		return Group{}, err
	}
	return group, nil
})

// JoinGroupAPI attempts to have the user with the given user id join the group with the given group code. It gets these values from the given GroupJoin struct.
// It returns the joined group if successful and an error if not
var JoinGroupAPI = NewAPI(http.MethodPost, "/api/joinGroup", func(groupJoin GroupJoin) (Group, error) {
	return JoinGroupDB(groupJoin)
})

// UpdateGroupAPI updates the name and members of the given group with the given values. It returns a confirmation string if successful and an error if not.
var UpdateGroupAPI = NewAPI(http.MethodPut, "/api/updateGroup", func(group Group) (string, error) {
	err := UpdateGroupDB(group)
	if err != nil {
		return "", err
	}
	return "group updated", nil
})

// GetMealsAPI returns the meals in the database that have the given group id.
var GetMealsAPI = NewAPI(http.MethodGet, "/api/getMeals", func(groupID int64) (Meals, error) {
	return GetMealsDB(groupID)
})

// CreateMealAPI creates and returns a new meal with the given information.
var CreateMealAPI = NewAPI(http.MethodPost, "/api/createMeal", func(meal Meal) (Meal, error) {
	return CreateMealDB(meal)
})

// UpdateMealAPI updates the given meal in the database to have the values of the given meal.
// It returns a confirmation string if successful and an error otherwise.
var UpdateMealAPI = NewAPI(http.MethodPut, "/api/updateMeal", func(meal Meal) (string, error) {
	err := UpdateMealDB(meal)
	if err != nil {
		return "", err
	}
	return "meal updated", nil
})

// DeleteMealAPI deletes the given meal from the database.
// It returns a confirmation string if successful and an error if not.
var DeleteMealAPI = NewAPI(http.MethodDelete, "/api/deleteMeal", func(id int64) (string, error) {
	err := DeleteMealEntriesDB(id)
	if err != nil {
		return "", err
	}
	err = DeleteMealDB(id)
	if err != nil {
		return "", err
	}
	return "meal deleted", nil
})

// GetEntriesAPI fetches and returns the entries associated with the given group id from the database
var GetEntriesAPI = NewAPI(http.MethodGet, "/api/getEntries", func(groupID int64) (Entries, error) {
	return GetEntriesDB(groupID)
})

// GetEntriesForMealAPI fetches and returns the entries associated with the given meal id from the database
var GetEntriesForMealAPI = NewAPI(http.MethodGet, "/api/getEntriesForMeal", func(mealID int64) (Entries, error) {
	return GetEntriesForMealDB(mealID)
})

// CreateEntryAPI creates and returns a new entry with the given entry's meal and user id values
var CreateEntryAPI = NewAPI(http.MethodPost, "/api/createEntry", func(entry Entry) (Entry, error) {
	return CreateEntryDB(entry)
})

// UpdateEntryAPI updates the given entry in the database to have the given information.
// It returns a confirmation string if successful and an error if not.
var UpdateEntryAPI = NewAPI(http.MethodPut, "/api/updateEntry", func(entry Entry) (string, error) {
	err := UpdateEntryDB(entry)
	if err != nil {
		return "", err
	}
	return "entry updated", nil
})

// DeleteEntryAPI deletes the entry with the given id from the database.
// It returns a confirmation string if successful and an error if not.
var DeleteEntryAPI = NewAPI(http.MethodDelete, "/api/deleteEntry", func(id int64) (string, error) {
	err := DeleteEntryDB(id)
	if err != nil {
		return "", err
	}
	return "entry deleted", nil
})

// RecommendRecipesAPI returns a list of recommended recipes based on the given word score map
var RecommendRecipesAPI = NewAPI(http.MethodGet, "/api/recommendRecipes", func(wordScoreMap map[string]Score) (Recipes, error) {
	return RecommendRecipes(wordScoreMap), nil
})

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

var SignUpAPI = NewAPI(http.MethodPost, "/api/signUp", func(user User) (User, error) {
	if user.Username == "" || user.Password == "" {
		return User{}, errors.New("username and password must not be empty")
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

var SignOutAPI = NewAPI(http.MethodPost, "/api/signOut", func(user User) (string, error) {
	sessionHash := sha256.Sum256([]byte(user.Session))
	sessionHashString := hex.EncodeToString(sessionHash[:])
	err := DeleteSessionDB(sessionHashString)
	if err != nil {
		return "", err
	}
	return "signed out", nil
})

var GetMealsAPI = NewAPI(http.MethodGet, "/api/getMeals", func(userID int) (Meals, error) {
	return GetMealsDB(userID)
})

var CreateMealAPI = NewAPI(http.MethodPost, "/api/createMeal", func(userID int) (Meal, error) {
	return CreateMealDB(userID)
})

var UpdateMealAPI = NewAPI(http.MethodPost, "/api/updateMeal", func(meal Meal) (string, error) {
	err := UpdateMealDB(meal)
	if err != nil {
		return "", err
	}
	return "meal updated", nil
})

var DeleteMealAPI = NewAPI(http.MethodDelete, "/api/deleteMeal", func(id int) (string, error) {
	err := DeleteMealDB(id)
	if err != nil {
		return "", err
	}
	return "meal deleted", nil
})

var GetPeopleAPI = NewAPI(http.MethodGet, "/api/getPeople", func(userID int) (People, error) {
	return GetPeopleDB(userID)
})

var CreatePersonAPI = NewAPI(http.MethodPost, "/api/createPerson", func(userID int) (Person, error) {
	return CreatePersonDB(userID)
})

var UpdatePersonAPI = NewAPI(http.MethodPost, "/api/updatePerson", func(person Person) (string, error) {
	err := UpdatePersonDB(person)
	if err != nil {
		return "", err
	}
	return "person updated", nil
})

var DeletePersonAPI = NewAPI(http.MethodDelete, "/api/deletePerson", func(id int) (string, error) {
	err := DeletePersonDB(id)
	if err != nil {
		return "", err
	}
	return "person deleted", nil
})

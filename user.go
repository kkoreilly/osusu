package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// User is a struct that represents a user in the users database
type User struct {
	ID       int
	Username string
	Password string
}

// CreateUserRequest sends an HTTP request to the server to create a user and returns the created user if successful and an error if not
func CreateUserRequest(user User) (User, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return User{}, err
	}
	req, err := NewRequest(http.MethodPost, "/api/createUser", bytes.NewBuffer(jsonData))
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	// return gotten user
	var res User
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return User{}, err
	}
	return res, nil
}

// SignInRequest sends an HTTP request to the server to sign in a user and returns the signed in user if successful and an error if not
func SignInRequest(user User) (User, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return User{}, err
	}
	req, err := NewRequest(http.MethodPost, "/api/signIn", bytes.NewBuffer(jsonData))
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	// return gotten user
	var res User
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return User{}, err
	}
	return res, nil
}

// AutoSignIn attempts to sign in using the saved login information, if it exists. Returns nil if signed in, an error if not.
func AutoSignIn(ctx app.Context) error {
	user := GetCurrentUser(ctx)
	if user.Username == "" || user.Password == "" {
		return errors.New("no saved login information")
	}
	user, err := SignInRequest(user)
	if err != nil {
		return err
	}
	SetCurrentUser(user, ctx)
	return nil
}

// SetCurrentUser sets the value of the current user in local storage
func SetCurrentUser(user User, ctx app.Context) {
	ctx.SetState("currentUser", user, app.Persist, app.Encrypt)
}

// GetCurrentUser gets the value of the current user from local storage
func GetCurrentUser(ctx app.Context) User {
	var user User
	ctx.GetState("currentUser", &user)
	return user
}

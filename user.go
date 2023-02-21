package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// User is a struct that represents a user in the users database
type User struct {
	ID         int
	Username   string
	Password   string
	Session    string // session id, not part of user in user database, but stored locally
	RememberMe bool   // also not part of user database, but used to transmit whether to save session
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
	// if no error, we are now authenticated
	authenticated = time.Now()
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
	// if no error, we are now authenticated
	authenticated = time.Now()
	return res, nil
}

// SignOutRequest sends an HTTP request to the server to sign out a user
func SignOutRequest(user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	req, err := NewRequest(http.MethodPost, "/api/signOut", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s: %v", resp.Status, string(body))
	}
	// if no error, we are no longer authenticated
	authenticated = time.UnixMilli(0)
	return nil
}

// SetCurrentUser sets the value of the current user in local storage
func SetCurrentUser(user User, ctx app.Context) {
	if user.RememberMe {
		ctx.SetState("currentUser", user, app.Persist, app.Encrypt, app.ExpiresIn(RememberMeSessionLength))
		return
	}
	ctx.SetState("currentUser", user, app.ExpiresIn(TemporarySessionLength))
}

// GetCurrentUser gets the value of the current user from local storage
func GetCurrentUser(ctx app.Context) User {
	var user User
	ctx.GetState("currentUser", &user)
	return user
}

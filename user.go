package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// User is a struct that represents a user in the users database
type User struct {
	ID       int
	Username string
	Password string
}

// CreateUserRequest sends an HTTP request to the server to create a user
func CreateUserRequest(user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	resp, err := http.Post("/api/createUser", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("create user request failed with status %s and body %v", resp.Status, string(body))
	}
	return nil
}

// SignInRequest sends an HTTP request to the server to sign in a user
func SignInRequest(user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	resp, err := http.Post("/api/signIn", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("sign in request failed with status %s and body %v", resp.Status, string(body))
	}
	return nil
}

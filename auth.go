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

// authenticated is when, if ever, the user has already been authenticated in this session of the app. This information is used to skip unnecessary additional authentication requests in the same session.
var authenticated time.Time

// Authenticate checks whether the user is signed in and takes an action or takes no action based on that. It returns whether the calling function should return.
// If required is set to true, Auth does nothing if the user is signed in and redirects the user to the sign in page otherwise.
// If required is set to false, Auth redirects the user to the home page if the user is signed in, and does nothing otherwise.
func Authenticate(required bool, ctx app.Context) bool {
	ok := time.Since(authenticated) < time.Hour
	if !ok {
		user := GetCurrentUser(ctx)
		if user.Session != "" {
			err := AuthenticateRequest(user)
			if err == nil {
				ok = true
				authenticated = time.Now()
			}

		}
	}
	switch {
	case required && !ok:
		ctx.Navigate("/signin")
	case !required && ok:
		ctx.Navigate("/home")
	default:
		return false
	}
	return true
}

// AuthenticateRequest sends an HTTP request to the server to authenticate a user based on its session id
func AuthenticateRequest(user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	req, err := NewRequest(http.MethodPost, "/api/authenticateSession", bytes.NewBuffer(jsonData))
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
	return nil
}

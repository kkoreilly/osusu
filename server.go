package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"golang.org/x/crypto/bcrypt"
)

func startServer() {
	http.Handle("/", &app.Handler{
		Name:        "MealRec",
		Title:       "MealRec",
		Icon:        app.Icon{Default: "/web/images/icon-192.png", Large: "/web/images/icon-512.png"},
		Description: "An app for getting recommendations on what meals to eat",
		Styles:      []string{"https://fonts.googleapis.com/css?family=Roboto", "/web/css/global.css", "/web/css/start.css", "/web/css/signinup.css", "/web/css/home.css", "/web/css/edit.css", "/web/css/people.css"},
	})

	err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/createUser", handleCreateUser)
	http.HandleFunc("/api/signIn", handleSignIn)
	http.HandleFunc("/api/getMeals", handleGetMeals)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("The %s method is not supported for this URL; only the POST method is supported for this URL.", r.Method), http.StatusBadRequest)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "username and password must not be empty", http.StatusBadRequest)
		return
	}
	originalPassword := user.Password
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Password = string(hash)

	user, err = CreateUserDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return original password instead of encrypted version
	user.Password = originalPassword
	// return user
	json, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("The %s method is not supported for this URL; only the POST method is supported for this URL.", r.Method), http.StatusBadRequest)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = SignInDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// return user
	json, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleGetMeals(w http.ResponseWriter, r *http.Request) {
	userIDParams := r.URL.Query()["u"]
	if userIDParams == nil {
		http.Error(w, "missing user id parameter", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDParams[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := User{ID: userID}

	meals, err := GetMealsDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(meals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

package main

import (
	"encoding/json"
	"io"
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

	handleFunc(http.MethodPost, "/api/createUser", handleCreateUser)
	handleFunc(http.MethodPost, "/api/signIn", handleSignIn)
	handleFunc(http.MethodGet, "/api/getMeals", handleGetMeals)
	handleFunc(http.MethodPost, "/api/createMeal", handleCreateMeal)
	handleFunc(http.MethodPost, "/api/updateMeal", handleUpdateMeal)
	handleFunc(http.MethodDelete, "/api/deleteMeal", handleDeleteMeal)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func handleFunc(method string, path string, handler http.HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, "The "+r.Method+" method is not supported for this URL; only the "+method+" method is supported for this URL.", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
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
	ownerParams := r.URL.Query()["u"]
	if ownerParams == nil {
		http.Error(w, "missing user id parameter", http.StatusBadRequest)
		return
	}
	owner, err := strconv.Atoi(ownerParams[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meals, err := GetMealsDB(owner)
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

func handleCreateMeal(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	owner, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meal, err := CreateMealDB(owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return meal
	json, err := json.Marshal(meal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleUpdateMeal(w http.ResponseWriter, r *http.Request) {
	var meal Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = UpdateMealDB(meal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("meal updated"))
}

func handleDeleteMeal(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = DeleteMealDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("meal deleted"))
}

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"golang.org/x/crypto/bcrypt"
)

// API access constants
const (
	APIUsername = "z3J8i6gVMyA!H$Ukpvqt5xLos5FgicTeWYf*MtfFU48HMUeMTaMCN59biD^3VxBup@^n7wnWgzCg442!95R9QHnt^6uKZ7f5ip2ycUjbfQ3sWzCZWVP8xgw!dZTn!trD"
	APIPassword = "gbx5T3*UJSALdxAES$n@w2m6b4o949XKMHsApk@Zt4&q3cf$37Jvf#g4#nd95hSnc4K%#h!JD9ifSkDhQyPMT@brtuU!cFxBJwny!ukC$s^ZVPdPzkJm8DvX4bK7to7d"
)

var sessions = make(map[string]int)

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

	HandleFunc(http.MethodPost, "/api/createUser", handleCreateUser)
	HandleFunc(http.MethodPost, "/api/signIn", handleSignIn)
	HandleFunc(http.MethodPost, "/api/authenticateSession", handleAuthenticateSession)

	HandleFunc(http.MethodGet, "/api/getMeals", handleGetMeals)
	HandleFunc(http.MethodPost, "/api/createMeal", handleCreateMeal)
	HandleFunc(http.MethodPost, "/api/updateMeal", handleUpdateMeal)
	HandleFunc(http.MethodDelete, "/api/deleteMeal", handleDeleteMeal)

	HandleFunc(http.MethodGet, "/api/getPeople", handleGetPeople)
	HandleFunc(http.MethodPost, "/api/createPerson", handleCreatePerson)
	HandleFunc(http.MethodPost, "/api/updatePerson", handleUpdatePerson)
	HandleFunc(http.MethodDelete, "/api/deletePerson", handleDeletePerson)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

// HandleFunc handles the given path with the given handler, requiring basic auth and accepting only the given method
func HandleFunc(method string, path string, handler http.HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(APIUsername))
		expectedPasswordHash := sha256.Sum256([]byte(APIPassword))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if !(usernameMatch && passwordMatch) {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, "The "+r.Method+" method is not supported for this URL; only the "+method+" method is supported for this URL.", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// GenerateSessionID generates a session id
func GenerateSessionID() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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
	// make session id
	session, err := GenerateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Session = session
	sessionHash := sha256.Sum256([]byte(session))
	err = CreateSessionDB(hex.EncodeToString(sessionHash[:]), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return blank password
	user.Password = ""
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
	// make session
	session, err := GenerateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	user.Session = session
	sessionHash := sha256.Sum256([]byte(session))
	err = CreateSessionDB(hex.EncodeToString(sessionHash[:]), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return blank password
	user.Password = ""
	// return user
	json, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleAuthenticateSession(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sessionHash := sha256.Sum256([]byte(user.Session))
	userID, expires, err := GetSessionDB(hex.EncodeToString(sessionHash[:]))
	if err != nil || userID != user.ID || expires.Before(time.Now()) {
		http.Error(w, "Invalid Session ID", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Authenticated"))
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

	meals, err := GetMealsDB(userID)
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
	userID, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	meal, err := CreateMealDB(userID)
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

func handleGetPeople(w http.ResponseWriter, r *http.Request) {
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

	people, err := GetPeopleDB(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleCreatePerson(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person, err := CreatePersonDB(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return person
	json, err := json.Marshal(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func handleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = UpdatePersonDB(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("person updated"))
}

func handleDeletePerson(w http.ResponseWriter, r *http.Request) {
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
	err = DeletePersonDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("person deleted"))
}

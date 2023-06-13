package server

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kkoreilly/osusu/db"
	"github.com/kkoreilly/osusu/osusu"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// AllRecipes are all of the recipes external recipes stored in the web/data/recipes.json file
var AllRecipes osusu.Recipes

func Start() {
	http.Handle("/", &app.Handler{
		Name:  "Osusu",
		Title: "Osusu",
		Icon: app.Icon{
			Default:    "/web/images/icon-192.png",
			Large:      "/web/images/icon-512.png",
			AppleTouch: "/web/images/icon-apple-touch.png",
		},
		Description: "Osusu is an app for getting recommendations on what meals to eat in a group based on the ratings of each member of the group, and the cost, effort, healthiness, and recency of the meal.",
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Roboto&family=Material+Symbols+Outlined",
			"/web/css/global.css",
			"/web/css/page.css",
			"/web/css/home.css",
			"/web/css/entries.css",
			"/web/css/entry.css",
			"/web/css/groups.css",
			"/web/css/group.css",
			"/web/css/input.css",
			"/web/css/button.css",
			"/web/css/chips.css",
			"/web/css/account.css",
			"/web/css/pie.css",
			"/web/css/start.css",
			"/web/css/recipe.css",
			"/web/css/mealimage.css",
		},
		AutoUpdateInterval: 10 * time.Second,
	})

	err := db.Start()
	if err != nil {
		log.Fatal(err)
	}

	osusu.InitRecipeConstants()
	recipes, err := osusu.LoadRecipes("web/data/recipes.json")
	if err != nil {
		log.Println(fmt.Errorf("error loading recipes: %w", err))
	}
	AllRecipes = recipes
	AllRecipes = AllRecipes.ComputeBaseScores()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
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
		// we hash the username and password so that they are guarenteed to be the same length, which is needed for subtle comparison algorithm
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		// username and password should already be hashed once on client side, so we have to hash the expected username and password twice to get the same value
		expectedUsernameHash := sha256.Sum256([]byte(osusu.APIUsername))
		expectedPasswordHash := sha256.Sum256([]byte(osusu.APIPassword))

		secondExpectedUsernameHash := sha256.Sum256(expectedUsernameHash[:])
		secondExpectedPasswordHash := sha256.Sum256(expectedPasswordHash[:])

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], secondExpectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], secondExpectedPasswordHash[:]) == 1)

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

// GenerateGroupCode generates a group code
func GenerateGroupCode() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

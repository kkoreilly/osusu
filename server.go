package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// API access constants
const (
	APIUsername = "z3J8i6gVMyA!H$Ukpvqt5xLos5FgicTeWYf*MtfFU48HMUeMTaMCN59biD^3VxBup@^n7wnWgzCg442!95R9QHnt^6uKZ7f5ip2ycUjbfQ3sWzCZWVP8xgw!dZTn!trD"
	APIPassword = "gbx5T3*UJSALdxAES$n@w2m6b4o949XKMHsApk@Zt4&q3cf$37Jvf#g4#nd95hSnc4K%#h!JD9ifSkDhQyPMT@brtuU!cFxBJwny!ukC$s^ZVPdPzkJm8DvX4bK7to7d"
)

func startServer() {
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
			"https://fonts.googleapis.com/css?family=Roboto",
			"/web/css/global.css",
			"/web/css/page.css",
			"/web/css/home.css",
			"/web/css/entries.css",
			"/web/css/entry.css",
			"/web/css/groups.css",
			"/web/css/group.css",
			"/web/css/input.css",
			"/web/css/chips.css",
			"/web/css/account.css",
		},
		AutoUpdateInterval: 10 * time.Second,
	})

	err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	CreateTablesDB()

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

// GenerateGroupCode generates a group code
func GenerateGroupCode() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

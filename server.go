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
	"strings"
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
		},
		AutoUpdateInterval: 10 * time.Second,
	})

	err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	err = InitDB()
	if err != nil {
		log.Fatal(err)
	}
	err = CleanupDB()
	if err != nil {
		log.Println(err)
	}

	recipes, err := GetRecipes()
	if err != nil {
		log.Println(err)
	} else {
		recipes = FixRecipeTimes(recipes)
		total := 0
		for i, recipe := range recipes {
			// if i > 100 {
			// 	break
			// }
			// fmt.Println("Index:", i, "Total Time:", recipe.TotalTime, "Prep Time:", recipe.PrepTime, "Cook Time:", recipe.CookTime, "Name:", recipe.Name)
			for _, cuisine := range []string{"African", "American", "Anglo-Indian", "Arabian", "Argentine", "Armenian", "Australian", "Austrian", "Azeri",
				"Balkan", "Bangladeshi", "Barbeque", "Basque", "Belgian", "Bengali", "Bhutanese", "Bolivian", "Brazilian", "British",
				"Bruneian", "Bulgarian", "Burmese", "Cambodian", "Cantonese", "Cape Malay", "Central Asian", "Cherokee", "Chilean",
				"Chinese", "Colombian", "Cornish", "Costa Rican", "Croatian", "Cuban", "Cypriot", "Czech", "Danish", "Djiboutian",
				"Dominican", "Dutch", "East African", "Eastern European", "Ecuadorian", "Egyptian", "Eritrean", "Estonian",
				"Ethiopian", "Faroe Islands", "Filipino", "Finnish", "French", "Galician", "Gambian", "Georgian", "German",
				"Ghanaian", "Greek", "Grenadian", "Guatemalan", "Guinea-Bissauan", "Guyanese", "Haitian", "Hawaiian", "Herzegovinian",
				"Hungarian", "Icelandic", "Indian", "Indonesian", "Iranian", "Iraqi", "Irish", "Israeli", "Italian",
				"Jamaican", "Japanese", "Jordanian", "Kazakh", "Kenyan", "Khmer", "Korean", "Kosovan", "Kuwaiti",
				"Kyrgyz", "Laotian", "Latin American", "Latvian", "Lebanese", "Lithuanian", "Luxembourgish", "Macedonian",
				"Malagasy", "Malaysian", "Maldivian", "Maltese", "Marshallese", "Mauritanian", "Mauritian", "Mexican",
				"Micronesian", "Middle Eastern", "Mongolian", "Moroccan", "Mozambican", "Myanmar", "Namibian", "Nepalese",
				"New Zealand", "Nicaraguan", "Nigerian", "North African", "North American", "Norwegian", "Omani", "Pakistani",
				"Palauan", "Palestinian", "Panamanian", "Papua New Guinean", "Paraguayan", "Peruvian", "Philippine",
				"Polish", "Portuguese", "Qatari", "Romanian", "Russian", "Rwandan", "Saint Lucian", "Salvadoran", "Samoa",
				"Samoan", "Sanmarinese", "Sao Tome and Principe", "Saudi Arabian", "Scottish", "Senegalese", "Serbian",
				"Seychellois", "Sierra Leonean", "Singaporean", "Slovak", "Slovenian", "Solomon Islander", "Somali",
				"South African", "South American", "South Korean", "Spanish", "Sri Lankan", "Sudanese", "Surinamese",
				"Swazi", "Swedish", "Swiss", "Syrian", "Taiwanese", "Tajikistani", "Tanzanian", "Thai", "Tibetan",
				"Tonga", "Trinidad and Tobago", "Tunisian", "Turkish", "Turkmen", "Tuvaluan", "Ugandan", "Ukrainian",
				"Uruguayan", "Uzbek", "Vietnamese", "Welsh", "West African", "Western European", "Yemeni", "Zambian",
				"Zimbabwean"} {
				if strings.Contains(recipe.Name, cuisine) || strings.Contains(recipe.Description, cuisine) || strings.Contains(recipe.Ingredients, cuisine) {
					total++
					log.Println(i, cuisine)
				}
			}
			// if recipe.PrepTime != "" {
			// 	totalTimeString := recipe.PrepTime[2:]
			// 	totalTimeString = strings.ToLower(totalTimeString)
			// 	totalTime, err := time.ParseDuration(totalTimeString)
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	log.Println("Index:", i, "Name:", recipe.Name, "Total Time:", totalTime)
			// }

		}
		log.Println(total)
	}

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

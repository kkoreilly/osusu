// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

//go:embed google-secret.json
var googleSecret []byte

// //go:embed google-icon.svg
// var googleIcon embed.FS

// func init() {
// 	icons.Icons = merged_fs.NewMergedFS(icons.Icons, googleIcon)
// }

/*
// Google authenticates the user with Google.
func Google() (*oauth2.Token, error) {
	ctx := context.TODO()

	config, err := google.ConfigFromJSON(googleSecret, "openid")
	if err != nil {
		return nil, err
	}

	// TODO(kai/auth): is this a good way to determine the port?
	port := rand.Intn(10_000)
	sport := ":" + strconv.Itoa(port)
	config.RedirectURL += sport

	code := make(chan string)
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code <- r.URL.Query().Get("code")
		w.Write([]byte("<h1>Signed in</h1><p>You can return to the app</p>"))
	})
	// TODO(kai/auth): more graceful closing / error handling
	go http.ListenAndServe(sport, sm)

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	goosi.TheApp.OpenURL(url)

	cs := <-code
	token, err := config.Exchange(ctx, cs, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, err
	}
	return token, nil
}
*/

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

func Google(ctx context.Context) (*oauth2.Token, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://127.0.0.1:5556/auth/google/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	code := make(chan string)

	sm := http.NewServeMux()
	sm.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		state, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}
		code <- r.URL.Query().Get("code")
		w.Write([]byte("<h1>Signed in</h1><p>You can return to the app</p>"))
	})
	// TODO(kai/auth): more graceful closing / error handling
	go http.ListenAndServe(":5556", sm)

	cs := <-code

	oauth2Token, err := config.Exchange(ctx, cs)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return oauth2Token, fmt.Errorf("failed to get user info: %w", err)
	}
	fmt.Println(userInfo)
	return oauth2Token, nil
}

// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"goki.dev/gi/v2/gi"
	"goki.dev/goosi"
	"goki.dev/goosi/events"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//go:embed google-secret.json
var googleSecret embed.FS

// GoogleButton adds a new button for signing in with Google.
func GoogleButton(par gi.Widget) *gi.Button {
	bt := gi.NewButton(par, "sign-in-with-google").SetType(gi.ButtonOutlined).SetText("Sign in with Google")
	bt.OnClick(func(e events.Event) {
		err := Google()
		if err != nil {
			gi.NewDialog(par).Title("Error signing in with Google").Prompt(err.Error())
		}
	})
	return bt
}

// Google authenticates the user with Google.
func Google() error {
	ctx := context.TODO()

	b, err := googleSecret.ReadFile("google-secret.json")
	if err != nil {
		return err
	}
	config, err := google.ConfigFromJSON(b, "openid")
	if err != nil {
		return err
	}

	// TODO(kai/auth): is this a good way to determine the port?
	port := rand.Intn(10_000)
	sport := ":" + strconv.Itoa(port)
	config.RedirectURL += sport

	code := make(chan string)
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code <- r.URL.Query().Get("code")
		w.Write([]byte("<h1>Authentication Successful</h1><p>You can close this browser tab and return to the app</p>"))
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
		return err
	}
	fmt.Println("got token", token)
	return nil
}

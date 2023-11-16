// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/goosi"
	"goki.dev/icons"
	"golang.org/x/oauth2"

	"github.com/yalue/merged_fs"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

//go:embed svg/google.svg
var googleIcon embed.FS

func init() {
	icons.Icons = merged_fs.NewMergedFS(icons.Icons, googleIcon)
}

// Google authenticates the user with Google and returns the oauth token
// and user info.
func Google(ctx context.Context) (*oauth2.Token, *oidc.UserInfo, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, nil, err
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://127.0.0.1:5556/auth/google/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	fmt.Println(config)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.RawURLEncoding.EncodeToString(b)

	code := make(chan string)

	sm := http.NewServeMux()
	sm.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}
		code <- r.URL.Query().Get("code")
		w.Write([]byte("<h1>Signed in</h1><p>You can return to the app</p>"))
	})
	// TODO(kai/auth): more graceful closing / error handling
	go http.ListenAndServe("127.0.0.1:5556", sm)

	goosi.TheApp.OpenURL(config.AuthCodeURL(state))

	cs := <-code

	oauth2Token, err := config.Exchange(ctx, cs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return oauth2Token, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return oauth2Token, userInfo, nil
}

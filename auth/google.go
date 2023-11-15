// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/goosi/events"
)

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
	return nil
}

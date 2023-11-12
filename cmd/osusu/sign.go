package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
)

type signInData struct {
	Username string
	Password string `view:"password"`
}

func signIn(sc *gi.Scene) {
	d := gi.NewDialog(sc, "sign-in").Title("Sign in")
	sd := &signInData{}
	giv.NewStructView(d).SetStruct(sd)
	d.Cancel().Ok("Sign in").Run()
}

type signUpData struct {
	Username string
	Password string `view:"password"`
}

func signUp(sc *gi.Scene) {
	d := gi.NewDialog(sc, "sign-up").Title("Sign up")
	sd := &signUpData{}
	giv.NewStructView(d).SetStruct(sd)
	d.Cancel().Ok("Sign up").Run()
}

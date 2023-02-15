package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Options represents meal recommendation options
type Options struct {
	CostWeight        int
	EffortWeight      int
	HealthinessWeight int
	TasteWeight       int
}

// SetOptions sets the options state value to the given meal recommendation options
func SetOptions(options Options, ctx app.Context) {
	ctx.SetState("options", options, app.Persist)
}

// GetOptions gets the meal recommendation options from local storage
func GetOptions(ctx app.Context) Options {
	var options Options
	ctx.GetState("options", &options)
	return options
}

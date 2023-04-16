package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

// Options represents meal recommendation options
type Options struct {
	People            map[int64]bool // key is person id, value is whether or not they are included
	Type              string
	Source            map[string]bool
	Cuisine           map[string]bool
	CostWeight        int
	EffortWeight      int
	HealthinessWeight int
	TasteWeight       int
	RecencyWeight     int
}

// DefaultOptions returns the default options for the given user
func DefaultOptions(user User) Options {
	options := Options{
		People:            make(map[int64]bool),
		Type:              "Dinner",
		Source:            map[string]bool{"Cooking": true, "Dine-In": true, "Takeout": true},
		Cuisine:           map[string]bool{},
		CostWeight:        50,
		EffortWeight:      50,
		HealthinessWeight: 50,
		TasteWeight:       50,
		RecencyWeight:     50,
	}
	for _, cuisine := range user.Cuisines {
		options.Cuisine[cuisine] = true
	}
	return options
}

// RemoveInvalidCuisines returns the the options with all invalid cuisines removed, using the given cuisine options
func (o Options) RemoveInvalidCuisines(cuisines []string) Options {
	res := map[string]bool{}
	for optionCuisine, value := range o.Cuisine {
		for _, cuisineOption := range cuisines {
			if value && optionCuisine == cuisineOption {
				res[optionCuisine] = true
			}
		}
	}
	o.Cuisine = res
	return o
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

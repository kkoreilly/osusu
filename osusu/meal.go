package osusu

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

//go:generate core generate -sql

type Meal struct {
	gorm.Model  `display:"-"`
	GroupID     uint  `display:"-"`
	Group       Group `display:"-"`
	Name        string
	Description string
	Image       string
	Source      Sources
	Category    Categories
	Cuisine     Cuisines
}

type Entry struct {
	gorm.Model  `display:"-"`
	MealID      uint `display:"-"`
	Meal        Meal `display:"-"`
	UserID      uint `display:"-"`
	User        User `display:"-"`
	Time        time.Time
	Category    Categories
	Source      Sources
	Taste       int `display:"slider" min:"0" def:"50" max:"100"`
	Cost        int `display:"slider" min:"0" def:"50" max:"100"`
	Effort      int `display:"slider" min:"0" def:"50" max:"100"`
	Healthiness int `display:"slider" min:"0" def:"50" max:"100"`
}

type Sources int64 //enums:bitflag

const (
	Cooking Sources = iota
	DineIn
	Takeout
	Delivery
)

type Categories int64 //enums:bitflag

const (
	Breakfast Categories = iota
	Brunch
	Lunch
	Dinner
	Dessert
	Snack
	Appetizer
	Side
	Drink
	Ingredient
)

type Cuisines int64 //enums:bitflag

const (
	African Cuisines = iota
	American
	Asian
	British
	Chinese
	European
	French
	Greek
	Indian
	Italian
	Japanese
	Jewish
	Korean
	LatinAmerican // Latin American
	Mexican
	MiddleEastern // Middle Eastern
	Thai
)

// Text returns all of the text associated with the meal as one string.
// It is intended to be used as text encoding model data, so it should
// not be presented to end-users.
func (m *Meal) Text() string {
	return strings.Join([]string{m.Name, m.Description}, "\n")
}

package osusu

import (
	"time"

	"gorm.io/gorm"
)

//go:generate enumgen -sql

type Meal struct {
	gorm.Model  `view:"-"`
	GroupID     uint  `view:"-"`
	Group       Group `view:"-"`
	Name        string
	Description string `view:"text-editor"`
	Image       string
	Source      Sources
	Category    Categories
	Cuisine     Cuisines
}

type Entry struct {
	gorm.Model  `view:"-"`
	MealID      uint `view:"-"`
	Meal        Meal `view:"-"`
	UserID      uint `view:"-"`
	User        User `view:"-"`
	Time        time.Time
	Category    Categories
	Source      Sources
	Taste       int `view:"slider" min:"0" def:"50" max:"100"`
	Cost        int `view:"slider" min:"0" def:"50" max:"100"`
	Effort      int `view:"slider" min:"0" def:"50" max:"100"`
	Healthiness int `view:"slider" min:"0" def:"50" max:"100"`
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
	LatinAmerican
	Mexican
	MiddleEastern
	Thai
)

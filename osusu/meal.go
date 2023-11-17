package osusu

import "gorm.io/gorm"

//go:generate enumgen -sql

type Meal struct {
	gorm.Model  `view:"-"`
	GroupID     uint  `view:"-"`
	Group       Group `view:"-"`
	Name        string
	Description string
	Source      string
	Image       string
	Category    Categories
	Cuisine     Cuisines
}

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

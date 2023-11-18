package osusu

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	MealID      uint
	Meal        Meal
	UserID      uint
	User        User
	Time        time.Time
	Category    Categories
	Source      Sources
	Cost        int
	Effort      int
	Healthiness int
	Taste       int
}

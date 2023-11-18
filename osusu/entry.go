package osusu

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model  `view:"-"`
	MealID      uint `view:"-"`
	Meal        Meal `view:"-"`
	UserID      uint `view:"-"`
	User        User `view:"-"`
	Time        time.Time
	Category    Categories
	Source      Sources
	Cost        int
	Effort      int
	Healthiness int
	Taste       int
}

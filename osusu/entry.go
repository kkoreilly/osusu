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
	Cost        int `view:"slider" min:"0" def:"50" max:"100"`
	Effort      int `view:"slider" min:"0" def:"50" max:"100"`
	Healthiness int `view:"slider" min:"0" def:"50" max:"100"`
	Taste       int `view:"slider" min:"0" def:"50" max:"100"`
}

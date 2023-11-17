package osusu

import "gorm.io/gorm"

type Meal struct {
	gorm.Model
	Name        string
	Description string
	Source      string
	Image       string
}

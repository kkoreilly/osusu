package osusu

import "gorm.io/gorm"

type Meal struct {
	gorm.Model  `view:"-"`
	Name        string
	Description string
	Source      string
	Image       string
}

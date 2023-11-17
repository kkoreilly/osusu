package osusu

import "gorm.io/gorm"

// User represents a user
type User struct {
	gorm.Model `view:"-"`
	Email      string
	Name       string
	Locale     string
	Picture    string
	GroupID    uint
}

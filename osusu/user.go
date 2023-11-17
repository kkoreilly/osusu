package osusu

import "gorm.io/gorm"

// User represents a user
type User struct {
	gorm.Model `view:"-"`
	GroupID    uint
	Email      string
	Name       string
	Locale     string
	Picture    string
}

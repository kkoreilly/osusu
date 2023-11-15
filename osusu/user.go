package osusu

import "gorm.io/gorm"

// User represents a user
type User struct {
	gorm.Model
	Username     string
	AccessToken  string
	RefreshToken string
}

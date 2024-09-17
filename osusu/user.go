package osusu

import "gorm.io/gorm"

type User struct {
	gorm.Model `display:"-"`
	GroupID    uint
	Email      string
	Name       string
	Locale     string
	Picture    string
}

/*
type Session struct {
	gorm.Model
	UserID uint
	User   User
	Token  string
}
*/

type Group struct {
	gorm.Model `display:"-"`
	Name       string
	Code       string `display:"-"`
	OwnerID    uint   `display:"-"`
	Owner      User   `display:"-"`
	Members    []User
}

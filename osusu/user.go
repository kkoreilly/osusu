package osusu

import "gorm.io/gorm"

type User struct {
	gorm.Model `view:"-"`
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
	gorm.Model `view:"-"`
	Name       string
	Code       string `view:"-"`
	OwnerID    uint   `view:"-"`
	Owner      User   `view:"-"`
	Members    []User
}

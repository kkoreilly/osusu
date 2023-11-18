package osusu

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserID uint
	User   User
	Token  string
}

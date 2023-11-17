package osusu

import "gorm.io/gorm"

type Group struct {
	gorm.Model `view:"-"`
	Name       string
	Code       string
	OwnerID    uint
	Owner      User
	Members    []User
}

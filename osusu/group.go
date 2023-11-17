package osusu

import "gorm.io/gorm"

type Group struct {
	gorm.Model `view:"-"`
	Name       string
	Code       string `view:"-"`
	OwnerID    uint   `view:"-"`
	Owner      User   `view:"-"`
	Members    []User
}

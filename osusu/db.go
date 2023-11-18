package osusu

import (
	"goki.dev/rqlite"
	"gorm.io/gorm"
)

// DB is the gorm database for the app.
var DB *gorm.DB

// OpenDB opens and sets up the database.
func OpenDB() error {
	db, err := gorm.Open(rqlite.Open("http://"))
	if err != nil {
		return err
	}
	DB = db
	return db.AutoMigrate(&User{}, &Session{}, &Group{}, &Meal{}, &Entry{})
}

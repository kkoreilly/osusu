// Package osusu contains the main code for the app that is common to both server and client code.
package osusu

import (
	"time"
)

// TemporarySessionLength is how long the user can be logged in without checking whether they are authenticated
const TemporarySessionLength = 3 * time.Hour

// RememberMeSessionLength is how long the user's session will be saved with remember me
const RememberMeSessionLength = 90 * 24 * time.Hour

package main

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"github.com/lib/pq/hstore"
)

// A UserMap holds information about what different users rate something.
// The key is the user id and the value is their rating.
type UserMap map[int64]int

// Sum returns the sum of the user map based on the given map specifying whether each of the users is participating and the given inversion.
// The key of the map is the user id and the value is whether they are participating.
// If inverted is true, the higher ratings of the users lead to a lower total sum.
func (u *UserMap) Sum(usersParticipating map[int64]bool, inverted bool) int {
	res := 0
	for i, v := range *u {
		participating := usersParticipating[i]
		// do not invert if the person is participating and we are not inverted or the person is not participating and we are inverted
		// otherwise, invert
		if participating && !inverted || !participating && inverted {
			res += v
		} else {
			res += 100 - v
		}
	}
	return res
}

// Average returns a simple average of the values of the user map
func (u *UserMap) Average() int {
	if len(*u) < 1 {
		return 0
	}
	res := 0
	for _, v := range *u {
		res += v
	}
	return res / len(*u)
}

// HasValueSet returns whether the given user has a value set for them in the user map
func (u *UserMap) HasValueSet(user User) bool {
	_, ok := (*u)[user.ID]
	return ok
}

// RemoveInvalid removes entries from the user map that are associated with nonexistent users
func (u *UserMap) RemoveInvalid(users []User) {
	for userID := range *u {
		got := false
		for _, user := range users {
			if userID == user.ID {
				got = true
			}
		}
		if !got {
			delete(*u, userID)
		}
	}
}

// Scan scans the provided value onto the user map
func (u *UserMap) Scan(value any) error {
	hStore := &hstore.Hstore{}
	err := hStore.Scan(value)
	if err != nil {
		return err
	}
	res := make(UserMap)
	for k, v := range hStore.Map {
		kInt, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil
		}
		vInt, err := strconv.Atoi(v.String)
		if err != nil {
			return nil
		}
		res[kInt] = vInt
	}
	*u = res
	return nil
}

// Value returns the database driver value of the given user map
func (u UserMap) Value() (driver.Value, error) {
	newMap := map[string]sql.NullString{}
	for k, v := range u {
		kString := strconv.FormatInt(k, 10)
		vString := strconv.Itoa(v)
		vNullString := sql.NullString{String: vString, Valid: vString != ""}
		newMap[kString] = vNullString
	}
	return hstore.Hstore{Map: newMap}.Value()
}

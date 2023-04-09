package main

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"github.com/lib/pq/hstore"
)

// A PersonMap holds information about what different people rate something.
// The key is the person id and the value is their rating.
type PersonMap map[int]int

// Sum returns the sum of the person map based on the given map specifying whether each of the people is participating and the given inversion.
// The key of the map is the person id and the value is whether they are participating.
// If inverted is true, the higher ratings of the people lead to a lower total sum.
func (p *PersonMap) Sum(peopleParticipating map[int]bool, inverted bool) int {
	res := 0
	for i, v := range *p {
		participating := peopleParticipating[i]
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

// Average returns a simple average of the values of the person map
func (p *PersonMap) Average() int {
	if len(*p) < 1 {
		return 0
	}
	res := 0
	for _, v := range *p {
		res += v
	}
	return res / len(*p)
}

// HasValueSet returns whether the given person has a value set for them in the person map
func (p *PersonMap) HasValueSet(person Person) bool {
	_, ok := (*p)[person.ID]
	return ok
}

// RemoveInvalid removes entries from the person map that are associated with nonexistent people
func (p *PersonMap) RemoveInvalid(people People) {
	for personID := range *p {
		got := false
		for _, person := range people {
			if personID == person.ID {
				got = true
			}
		}
		if !got {
			delete(*p, personID)
		}
	}
}

// Scan scans the provided value onto the PersonMap
func (p *PersonMap) Scan(value any) error {
	hStore := &hstore.Hstore{}
	err := hStore.Scan(value)
	if err != nil {
		return err
	}
	res := make(PersonMap)
	for k, v := range hStore.Map {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		vInt, err := strconv.Atoi(v.String)
		if err != nil {
			return nil
		}
		res[kInt] = vInt
	}
	*p = res
	return nil
}

// Value returns the database driver value of the given person map
func (p PersonMap) Value() (driver.Value, error) {
	newMap := map[string]sql.NullString{}
	for k, v := range p {
		kString := strconv.Itoa(k)
		vString := strconv.Itoa(v)
		vNullString := sql.NullString{String: vString, Valid: vString != ""}
		newMap[kString] = vNullString
	}
	return hstore.Hstore{Map: newMap}.Value()
}

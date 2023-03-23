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

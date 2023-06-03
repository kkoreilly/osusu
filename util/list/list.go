// Package list provides functions for generating formatted strings of slice and map lists
package list

import (
	"sort"
	"strconv"
)

// Slice returns a formatted string of the given slice list of items
func Slice(list []string) string {
	res := ""
	lenList := len(list)
	for i, l := range list {
		res += l
		if lenList != 2 && i != lenList-1 {
			res += ", "
		}
		if lenList == 2 && i == lenList-2 {
			res += " and "
		}
		if lenList > 2 && i == lenList-2 {
			res += "and "
		}
	}
	return res
}

// SliceNum returns a formatted string of the given slice list of items, replacing any number of the elements above the specified number with "n more"
func SliceNum(list []string, num int) string {
	if len(list) <= num {
		return Slice(list)
	}
	return Slice(append(list[:num-1], strconv.Itoa(len(list)-num+1)+" more"))
}

// Map returns a formatted string of the given map list of items in which the key is the item and the value is whether it should be included in the string
func Map(list map[string]bool) string {
	slice := []string{}
	for k, v := range list {
		if v {
			slice = append(slice, k)
		}
	}
	// sort to prevent constant switching
	sort.Strings(slice)
	return Slice(slice)
}

// MapNum returns a formatted string of the given list of items in which the key is the item and the value is whether it should be included in the string.
// MapNum limits the number of items to the provided number and adds "and n more" to the end if this limit is exceeded.
func MapNum(list map[string]bool, num int) string {
	slice := []string{}
	for k, v := range list {
		if v {
			slice = append(slice, k)
		}
	}
	// sort to prevent constant switching
	sort.Strings(slice)
	return SliceNum(slice, num)
}

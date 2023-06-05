// Package file provides simple functions for interacting with the filesystem.
package file

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
)

// Exists returns whether the file with the given path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

// RenameSafe renames (moves) oldpath to newpath. If newpath already exists, it finds the path newpath + n (ex: newpath3) with the lowest n that is not taken and renames oldpath to that.
func RenameSafe(oldpath, newpath string) error {
	if !Exists(newpath) {
		return os.Rename(oldpath, newpath)
	}
	for i := 0; ; i++ {
		intPath := newpath + strconv.Itoa(i)
		if !Exists(intPath) {
			return os.Rename(oldpath, intPath)
		}
	}
}

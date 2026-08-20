//go:build !windows

package main

import "errors"

func executableDirectory() (string, error) {
	return "", errors.New("path commands are only supported on Windows")
}

func updateUserPath(directory string, action pathAction) (bool, error) {
	return false, errors.New("path commands are only supported on Windows")
}

func reloadPath() error {
	return errors.New("path commands are only supported on Windows")
}

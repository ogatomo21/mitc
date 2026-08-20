package main

import (
	"errors"
	"io"
)

type pathAction uint8

const (
	pathActionNone pathAction = iota
	pathActionAdd
	pathActionRemove
)

func runPathCommand(action pathAction, out io.Writer) error {
	directory, err := executableDirectory()
	if err != nil {
		return err
	}

	changed, err := updateUserPath(directory, action)
	if err != nil {
		return err
	}
	if !changed {
		if action == pathActionAdd {
			_, _ = io.WriteString(out, "mitc's directory is already in the Windows user PATH.\n")
		} else {
			_, _ = io.WriteString(out, "mitc's directory was not in the Windows user PATH.\n")
		}
		return nil
	}
	if err := reloadPath(); err != nil {
		return err
	}
	if action == pathActionAdd {
		_, _ = io.WriteString(out, "Added mitc's directory to the Windows user PATH and reloaded it.\n")
	} else {
		_, _ = io.WriteString(out, "Removed mitc's directory from the Windows user PATH and reloaded it.\n")
	}
	return nil
}

func validatePathAction(action pathAction) error {
	if action != pathActionAdd && action != pathActionRemove {
		return errors.New("invalid PATH action")
	}
	return nil
}

package launchd

import (
	"errors"
)

type ExecError struct {
	UnderlyingError error
	Command         []string
	Stderr          string
	Stdout          string
}

func (err *ExecError) Error() string {
	return err.UnderlyingError.Error()
}

func (err *ExecError) Is(err2 error) bool {
	if err2.Error() == err.Stderr {
		return true
	}

	return errors.Is(err.UnderlyingError, err)
}

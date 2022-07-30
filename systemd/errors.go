package systemd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ExecError struct {
	UnderlyingError error
	Command         []string
	Stderr          string
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

func (err *ExecError) Unwrap() error {
	return err.UnderlyingError
}

func (err *ExecError) IsBetterWithSudo() bool {
	return strings.HasSuffix(err.Stderr, `\nTry re-running the command as root for richer errors.\n`)
}

func (err *ExecError) MatchLaunchdReason(reason string) bool {
	return strings.HasSuffix(strings.SplitN(err.Stderr, "\n", 2)[0], fmt.Sprintf(": %s", reason))
}

func (err *ExecError) MatchBadRequestRegex(reason *regexp.Regexp) bool {
	lines := strings.SplitN(err.Stderr, "\n", 2)

	if lines[0] != "Bad request." {
		return false
	}

	return reason.MatchString(lines[1])
}

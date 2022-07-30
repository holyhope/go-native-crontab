package systemd

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"strconv"

	"github.com/holyhope/god"
)

func Name(opts god.Options) (string, error) {
	if !opts.HasName() {
		return "", god.NewMissingOptionError("Name")
	}

	value := opts.Name()

	parent := path.Dir(value)
	if !fs.ValidPath(value) || len(value) == 0 || value == "." || value == ".." || parent != "." {
		return "", &god.InvalidOptionError{
			Key:   "Name",
			Value: value,
		}
	}

	return value, nil
}

func Scope(opts god.Options) (god.Scope, error) {
	if !opts.HasScope() {
		return god.Scope(0), god.NewMissingOptionError("Scope")
	}

	scope := opts.Scope()

	switch scope {
	case god.ScopeUser, god.ScopeSystem:
		return scope, nil
	default:
		return god.Scope(0), &god.InvalidOptionError{
			Key:   "Scope",
			Value: scope,
		}
	}
}

func Domain(scope god.Scope) (string, error) {
	switch scope {
	case god.ScopeSystem:
		return "system", nil
	case god.ScopeUser:
		return fmt.Sprintf("user/%d", os.Getuid()), nil
	default:
		return "", &god.InvalidOptionError{
			Key:   "Scope",
			Value: scope,
		}
	}
}

func Path(scope god.Scope, filename string) (string, error) {
	switch scope {
	case god.ScopeUser:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home: %w", err)
		}

		return path.Join(home, "Library/LaunchAgents", filename), nil
	case god.ScopeSystem:
		return path.Join("/Library/LaunchAgents", filename), nil
	default:
		return "", &god.InvalidOptionError{
			Key:   "Scope",
			Value: scope,
		}
	}
}

func Username(opts god.Options) (string, error) {
	u, err := user.LookupId(strconv.Itoa(opts.UserOwner()))
	if err != nil {
		return "", &god.InvalidOptionError{
			Key:   "UserOwner",
			Value: opts.UserOwner(),
		}
	}

	return u.Username, nil
}

// https://developer.apple.com/library/archive/technotes/tn2083/_index.html#//apple_ref/doc/uid/DTS10003794-CH1-SUBSUBSECTION5
func LimitLoadToSessionType(sessionType god.DarwinLimitLoadToSessionType) (string, error) {
	switch sessionType {
	case god.DarwinLimitLoadToSessionAqua:
		return "Aqua", nil
	case god.DarwinLimitLoadToSessionLoginWindow:
		return "LoginWindow", nil
	case god.DarwinLimitLoadToSessionBackground:
		return "Background", nil
	case god.DarwinLimitLoadToSessionStandardIO:
		return "StandardIO", nil
	default:
		return "", &god.InvalidOptionError{
			Key:   "DarwinLimitLoadToSessionType",
			Value: sessionType,
		}
	}
}

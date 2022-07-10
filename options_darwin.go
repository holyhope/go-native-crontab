//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package god

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"
)

func (name Name) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.name = string(name)
	}

	return nil
}

func (ut Type) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.unitType = string(ut)
	}

	return nil
}

func (state State) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.state = string(state)
	}

	return nil
}

func (description Description) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.description = string(description)
	}

	return nil
}

func (interval Interval) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.interval = time.Duration(interval)
	}

	return nil
}

// Command is the command of the unit.
func Command(cmd string, args ...string) FactoryOpts {
	return FactoryOptsFn(func(ctx context.Context, u Unit) error {
		if u, ok := u.(*unit); ok {
			u.command = append([]string{cmd}, args...)
			return nil
		}

		return nil
	})
}

func (envs Envs) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.envs = envs
	}

	return nil
}

func (scope Scope) Path(filename string) (string, error) {
	switch scope {
	case ScopeUser:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home: %w", err)
		}

		return path.Join(home, "Library/LaunchAgents", filename), nil
	case ScopeSystem:
		return path.Join("/Library/LaunchAgents", filename+".plist"), nil
	default:
		return "", UnkownScopeError(scope)
	}
}

func (scope Scope) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.scope = scope
	}

	return nil
}

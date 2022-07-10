//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package god

import (
	"context"
	"fmt"
	"os"
	"path"
)

func (name UnitName) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.name = string(name)
	}

	return nil
}

func (ut UnitType) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.unitType = string(ut)
	}

	return nil
}

func (state UnitState) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.state = string(state)
	}

	return nil
}

func (description UnitDescription) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.description = string(description)
	}

	return nil
}

// UnitCommand is the command of the unit.
func UnitCommand(cmd string, args ...string) FactoryOpts {
	return FactoryOptsFn(func(ctx context.Context, u Unit) error {
		if u, ok := u.(*unit); ok {
			u.command = append([]string{cmd}, args...)
			return nil
		}

		return nil
	})
}

func (envs UnitEnvs) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.envs = envs
	}

	return nil
}

func (scope UnitScope) Path(filename string) (string, error) {
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

func (scope UnitScope) Apply(ctx context.Context, u Unit) error {
	if u, ok := u.(*unit); ok {
		u.scope = scope
	}

	return nil
}

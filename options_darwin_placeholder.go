//go:build !darwin && linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build !darwin,linux,!freebsd,!netbsd,!openbsd,!windows,!js

package god

import (
	"context"
)

func (s DarwinOptsScope) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	return ct, nil
}

func (f FileName) Apply(ctx context.Context, ct CronTab) (CronTab, error) {
	return ct, nil
}

// UnitCommand is the command of the unit.
func UnitCommand(cmd string, args ...string) FactoryOpts {
	return func(ctx context.Context, u Unit) (Unit, error) {
		return u, nil
	}
}

//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js
// +build darwin,!linux,!freebsd,!netbsd,!openbsd,!windows,!js

package launchd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/holyhope/god"
)

func (u *Unit) Install(ctx context.Context) error {
	if err := u.writeUnitFile(u.unitPath); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}

	if _, err := u.exec(ctx, u.launchctlPath, "bootstrap", u.domain, u.unitPath); err != nil {
		if err, ok := err.(*ExecError); ok {
			if err.MatchLaunchdReason("service already bootstrapped") {
				return nil
			}
		}

		return err
	}

	return nil
}

func (u *Unit) Uninstall(ctx context.Context) error {
	_, err := u.exec(ctx, u.launchctlPath, "bootout", u.domain, u.unitPath)
	if err != nil {
		if err, ok := err.(*ExecError); ok {
			if err.MatchLaunchdReason("No such file or directory") {
				return nil
			}
		}

		return err
	}

	if err := os.Remove(u.unitPath); err != nil {
		return fmt.Errorf("remove unit file: %w", err)
	}

	return nil
}

func (u *Unit) Reload(ctx context.Context) error {
	if err := u.Uninstall(ctx); err != nil {
		return err
	}
	if err := u.Install(ctx); err != nil {
		return err
	}

	return nil
}

// Status returns the status of the unit.
// Carefull, this rely on launchctl output which is not stable API yet.
func (u *Unit) Status(ctx context.Context) (god.UnitStatus, error) {
	_, err := u.exec(ctx, u.launchctlPath, "print", fmt.Sprintf("%s/%s", u.domain, u.unitName))
	if err != nil {
		if err, ok := err.(*ExecError); ok {
			if err.MatchBadRequestRegex(regexp.MustCompile(`Could not find service ".*" in domain .*`)) {
				if _, err := os.Stat(u.unitPath); err != nil {
					if os.IsNotExist(err) {
						return god.NewUnitStatus(false, false), nil
					}

					return nil, fmt.Errorf("stat unit file: %w", err)
				}

				return god.NewUnitStatus(true, false), nil
			}
		}

		return nil, fmt.Errorf("print unit info: %w", err)
	}

	return god.NewUnitStatus(true, true), nil
}

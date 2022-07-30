package systemd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	systemdBus "github.com/coreos/go-systemd/v22/dbus"
	"github.com/holyhope/god"
)

func (u *Unit) Create(ctx context.Context) error {
	if err := u.writeUnitFile(u.unitPath); err != nil {
		return err
	}

	return nil
}

func (u *Unit) Delete(ctx context.Context) error {
	if err := u.Disable(ctx); err != nil {
		return err
	}

	return nil
}

func (u *Unit) Enable(ctx context.Context) error {
	conn, err := systemdBus.NewWithContext(ctx)
	if err != nil {
		return fmt.Errorf("new systemd bus: %w", err)
	}

	defer conn.Close()

	ok, _, err := conn.EnableUnitFilesContext(ctx, []string{u.unitPath}, false, false)
	if err != nil {
		return fmt.Errorf("enable %s: %w", u.unitPath, err)
	}

	if !ok {
		return fmt.Errorf("%s not enabled: %w", u.unitName, err)
	}

	/*
		if _, err := systemdDaemon.SdNotify(true, systemdDaemon.SdNotifyReady); err != nil {
			return fmt.Errorf("sd-notify: %w", err)
		}
	*/

	return nil
}

func (u *Unit) Disable(ctx context.Context) error {
	conn, err := systemdBus.NewWithContext(ctx)
	if err != nil {
		return fmt.Errorf("new systemd bus: %w", err)
	}

	defer conn.Close()

	if _, err := conn.DisableUnitFilesContext(ctx, []string{u.unitPath}, false); err != nil {
		return fmt.Errorf("disable %s: %w", u.unitPath, err)
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
